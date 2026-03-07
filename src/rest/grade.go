package rest

import (
	"fmt"
	"math"
	"skulla-api/db"
	"skulla-api/services"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type SubjectGrade struct {
	Name             string  `json:"Name"`
	Marks            float64 `json:"Marks"`
	EvaluationsDone  int     `json:"EvaluationsDone"`
	EvaluationsTotal int     `json:"EvaluationsTotal"`
}

type GradesSummary struct {
	Subjects       []SubjectGrade `json:"Subjects"`
	FrequencyGrade float64        `json:"FrequencyGrade"`
	ExamGrade      float64        `json:"ExamGrade"`
	FinalGrade     *float64       `json:"FinalGrade"`
}

type StudentInfo struct {
	ID        uint   `json:"ID"`
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
}

type RegistrationGradeReport struct {
	ID         uint          `json:"ID"`
	Status     string        `json:"Status"`
	EnrolledAt *time.Time    `json:"EnrolledAt"`
	Student    StudentInfo   `json:"Student"`
	Grades     GradesSummary `json:"Grades"`
}

func ListGrades(c *fiber.Ctx) error {
	studentClassId, err := ParseUintQueryParam(c, "student_class_id", true)
	if err != nil {
		return ReturnBadRequest(c, err.Error())
	}

	userEmail, err := GetUserEmailFromToken(c)
	if err != nil {
		return ReturnBadRequest(c, err.Error())
	}

	courseID, err := db.GetStudentClassCourseID(studentClassId)
	if err != nil {
		return ReturnBadRequest(c, "Student class not found")
	}

	if !db.IsTeacherEmailBelongToCourse(userEmail, int(courseID)) {
		return ReturnUnauthorized(c, "User does not have permission to access course")
	}

	registrations := db.ListRegistrations(int(studentClassId))

	var studentClass db.StudentClass
	if err := db.DB().First(&studentClass, studentClassId).Error; err != nil {
		return ReturnBadRequest(c, "Student class not found")
	}

	var formula db.GradingFormula
	formulaQuery := db.DB().Where("course_id = ? AND is_active = ?", courseID, true)
	if studentClass.LevelId > 0 {
		formulaQuery = formulaQuery.Where("level_id = ?", studentClass.LevelId)
	} else {
		formulaQuery = formulaQuery.Where("level_id IS NULL")
	}
	formulaQuery.First(&formula)

	formulaCategoryNames := extractFormulaCategoryNames(formula, studentClass.LevelId)
	examCategoryNames := extractExamCategoryNames(formula, studentClass.LevelId)

	allCategories := db.FindActiveCategories(courseID, studentClass.LevelId)

	var categories []db.EvaluationCategory
	if len(formulaCategoryNames) > 0 {
		for _, cat := range allCategories {
			if formulaCategoryNames[cat.Name] {
				categories = append(categories, cat)
			}
		}
	} else {
		categories = allCategories
	}

	examCategoryIDs := make(map[uint]bool)
	for _, cat := range allCategories {
		if examCategoryNames[cat.Name] {
			examCategoryIDs[cat.ID] = true
		}
	}

	var allGrades []db.Grade
	db.DB().Where("student_class_id = ?", studentClassId).Find(&allGrades)

	gradesByRegAndCat := make(map[uint]map[uint][]db.Grade)
	for _, grade := range allGrades {
		if gradesByRegAndCat[grade.RegistrationID] == nil {
			gradesByRegAndCat[grade.RegistrationID] = make(map[uint][]db.Grade)
		}
		gradesByRegAndCat[grade.RegistrationID][grade.CategoryID] = append(
			gradesByRegAndCat[grade.RegistrationID][grade.CategoryID], grade)
	}

	calculator := services.NewGradingCalculator()
	var results []RegistrationGradeReport

	for _, reg := range registrations {
		report := RegistrationGradeReport{
			ID:         reg.ID,
			Status:     reg.Status,
			EnrolledAt: reg.EnrolledAt,
			Student: StudentInfo{
				ID:        reg.Student.ID,
				FirstName: reg.Student.FirstName,
				LastName:  reg.Student.LastName,
			},
		}

		var subjects []SubjectGrade
		for _, cat := range categories {
			grades := gradesByRegAndCat[reg.ID][cat.ID]

			evaluationsDone := 0
			var totalScore float64
			for _, g := range grades {
				if g.Score != nil && !g.IsExcused {
					evaluationsDone++
					totalScore += *g.Score
				}
			}

			var marks float64
			if evaluationsDone > 0 {
				marks = math.Round(totalScore/float64(evaluationsDone)*100) / 100
			}

			subjects = append(subjects, SubjectGrade{
				Name:             cat.Name,
				Marks:            marks,
				EvaluationsDone:  evaluationsDone,
				EvaluationsTotal: cat.Cardinality,
			})
		}

		report.Grades.Subjects = subjects

		hasExamGrade := len(examCategoryIDs) == 0
		if !hasExamGrade {
			for catID := range examCategoryIDs {
				for _, g := range gradesByRegAndCat[reg.ID][catID] {
					if g.Score != nil {
						hasExamGrade = true
						break
					}
				}
				if hasExamGrade {
					break
				}
			}
		}

		formulaResult, err := calculator.CalculateFormulaDetailed(reg.ID, studentClassId, courseID)
		if err == nil && formulaResult != nil {
			if hasExamGrade {
				report.Grades.FinalGrade = &formulaResult.FinalScore
			}

			for _, stage := range formulaResult.Stages {
				if stage.Name == "MF" {
					for _, comp := range stage.Components {
						if comp.Name != "CF" && comp.Name != "MAC" {
							report.Grades.ExamGrade = comp.Score
						}
					}
				} else {
					if report.Grades.FrequencyGrade == 0 {
						report.Grades.FrequencyGrade = stage.Score
					}
				}
			}
		}

		results = append(results, report)
	}

	return c.JSON(results)
}

func extractFormulaCategoryNames(formula db.GradingFormula, levelId int) map[string]bool {
	names := make(map[string]bool)
	if formula.FormulaConfig == nil {
		return names
	}

	config := formula.FormulaConfig

	formulaType, _ := config["type"].(string)

	var stagesRaw interface{}
	switch formulaType {
	case "multi_stage_level_based":
		levelsRaw, ok := config["levels"].(map[string]interface{})
		if !ok {
			return names
		}
		levelKey := fmt.Sprintf("%d", levelId)
		levelConfig, ok := levelsRaw[levelKey].(map[string]interface{})
		if !ok {
			return names
		}
		stagesRaw = levelConfig["stages"]
	case "multi_stage":
		stagesRaw = config["stages"]
	case "weighted_average":
		categoriesRaw, ok := config["categories"].(map[string]interface{})
		if ok {
			for name := range categoriesRaw {
				names[name] = true
			}
		}
		return names
	default:
		return names
	}

	stages, ok := stagesRaw.(map[string]interface{})
	if !ok {
		return names
	}

	stageNames := make(map[string]bool)
	for stageName := range stages {
		stageNames[stageName] = true
	}

	for _, stageDataRaw := range stages {
		stageData, ok := stageDataRaw.(map[string]interface{})
		if !ok {
			continue
		}
		stageType, _ := stageData["type"].(string)
		switch stageType {
		case "weighted_average":
			categoriesRaw, ok := stageData["categories"].(map[string]interface{})
			if ok {
				for name := range categoriesRaw {
					names[name] = true
				}
			}
		case "simple_average":
			componentsRaw, _ := stageData["components"].([]interface{})
			for _, compRaw := range componentsRaw {
				compName, _ := compRaw.(string)
				if !stageNames[compName] {
					names[compName] = true
				}
			}
		}
	}

	return names
}

type EvaluationItem struct {
	ID    int      `json:"id"`
	Marks *float64 `json:"marks"`
	Date  *string  `json:"date"`
}

type CategoryGrades struct {
	ID          uint             `json:"id"`
	Name        string           `json:"name"`
	Evaluations []EvaluationItem `json:"evaluations"`
}

type ExamInfo struct {
	ID    uint     `json:"id"`
	Marks *float64 `json:"marks"`
	Date  *string  `json:"date"`
}

type RegistrationGradesResponse struct {
	Categories     []CategoryGrades `json:"categories"`
	Exam           *ExamInfo        `json:"exam"`
	FrequencyGrade float64          `json:"frequencyGrade"`
	FinalGrade     *float64         `json:"finalGrade"`
}

func GetGradesOfRegistration(c *fiber.Ctx) error {
	parsed, err := strconv.ParseUint(c.Params("registrationId"), 10, 64)
	if err != nil {
		return ReturnBadRequest(c, "invalid registrationId")
	}
	registrationId := uint(parsed)

	var registration db.Registration
	if err := db.DB().Preload("StudentClass").First(&registration, registrationId).Error; err != nil {
		return ReturnNotFound(c, "registration not found")
	}

	studentClassID := registration.StudentClassID
	courseID := registration.StudentClass.CourseID
	levelID := registration.StudentClass.LevelId

	var formula db.GradingFormula
	formulaQuery := db.DB().Where("course_id = ? AND is_active = ?", courseID, true)
	if levelID > 0 {
		formulaQuery = formulaQuery.Where("level_id = ?", levelID)
	} else {
		formulaQuery = formulaQuery.Where("level_id IS NULL")
	}
	formulaQuery.First(&formula)

	examCategoryNames := extractExamCategoryNames(formula, levelID)
	formulaCategoryNames := extractFormulaCategoryNames(formula, levelID)

	allCategories := db.FindActiveCategories(courseID, levelID)

	var categories []db.EvaluationCategory
	if len(formulaCategoryNames) > 0 {
		for _, cat := range allCategories {
			if formulaCategoryNames[cat.Name] {
				categories = append(categories, cat)
			}
		}
	} else {
		categories = allCategories
	}

	var grades []db.Grade
	db.DB().Where("student_class_id = ? AND registration_id = ?", studentClassID, registrationId).
		Order("created_at ASC").Find(&grades)

	gradesByCat := make(map[uint][]db.Grade)
	for _, g := range grades {
		gradesByCat[g.CategoryID] = append(gradesByCat[g.CategoryID], g)
	}

	var response RegistrationGradesResponse
	var examInfo *ExamInfo

	for _, cat := range categories {
		if examCategoryNames[cat.Name] {
			catGrades := gradesByCat[cat.ID]
			exam := &ExamInfo{ID: cat.ID}
			if len(catGrades) > 0 && catGrades[0].Score != nil {
				exam.Marks = catGrades[0].Score
				if catGrades[0].Date != nil {
					dateStr := catGrades[0].Date.Format("2006-01-02")
					exam.Date = &dateStr
				}
			}
			examInfo = exam
			continue
		}

		catGrades := gradesByCat[cat.ID]
		var evaluations []EvaluationItem
		for i := 0; i < cat.Cardinality; i++ {
			eval := EvaluationItem{ID: 0}
			if i < len(catGrades) && catGrades[i].Score != nil && !catGrades[i].IsExcused {
				eval.Marks = catGrades[i].Score
				eval.ID = int(catGrades[i].ID)
				if catGrades[i].Date != nil {
					dateStr := catGrades[i].Date.Format("2006-01-02")
					eval.Date = &dateStr
				}
			}
			evaluations = append(evaluations, eval)
		}
		response.Categories = append(response.Categories, CategoryGrades{
			ID:          cat.ID,
			Name:        cat.Name,
			Evaluations: evaluations,
		})
	}

	response.Exam = examInfo

	hasExamGrade := len(examCategoryNames) == 0 || (examInfo != nil && examInfo.Marks != nil)

	calculator := services.NewGradingCalculator()
	formulaResult, calcErr := calculator.CalculateFormulaDetailed(registrationId, studentClassID, courseID)
	if calcErr == nil && formulaResult != nil {
		if hasExamGrade {
			response.FinalGrade = &formulaResult.FinalScore
		}

		for _, stage := range formulaResult.Stages {
			if stage.Name == "MF" {
				continue
			}
			if response.FrequencyGrade == 0 {
				response.FrequencyGrade = stage.Score
			}
		}
	}

	return c.JSON(response)
}

func extractExamCategoryNames(formula db.GradingFormula, levelId int) map[string]bool {
	names := make(map[string]bool)
	if formula.FormulaConfig == nil {
		return names
	}

	config := formula.FormulaConfig
	formulaType, _ := config["type"].(string)

	var stagesRaw interface{}
	switch formulaType {
	case "multi_stage_level_based":
		levelsRaw, ok := config["levels"].(map[string]interface{})
		if !ok {
			return names
		}
		levelKey := fmt.Sprintf("%d", levelId)
		levelConfig, ok := levelsRaw[levelKey].(map[string]interface{})
		if !ok {
			return names
		}
		stagesRaw = levelConfig["stages"]
	case "multi_stage":
		stagesRaw = config["stages"]
	default:
		return names
	}

	stages, ok := stagesRaw.(map[string]interface{})
	if !ok {
		return names
	}

	stageNames := make(map[string]bool)
	for stageName := range stages {
		stageNames[stageName] = true
	}

	mfRaw, ok := stages["MF"]
	if !ok {
		return names
	}
	mfData, ok := mfRaw.(map[string]interface{})
	if !ok {
		return names
	}

	componentsRaw, _ := mfData["components"].([]interface{})
	for _, compRaw := range componentsRaw {
		compName, _ := compRaw.(string)
		if !stageNames[compName] {
			names[compName] = true
		}
	}

	return names
}

func CreateGrade(c *fiber.Ctx) error {
	parsed, err := strconv.ParseUint(c.Params("registrationId"), 10, 64)
	if err != nil {
		return ReturnBadRequest(c, "invalid registrationId")
	}
	registrationId := uint(parsed)

	var body struct {
		EvaluationCategoryID uint    `json:"evaluation_category_id"`
		Marks                float64 `json:"marks"`
		Date                 string  `json:"date"`
	}

	if err := c.BodyParser(&body); err != nil {
		return ReturnBadRequest(c, "invalid request body")
	}

	if body.EvaluationCategoryID == 0 {
		return ReturnBadRequest(c, "evaluation_category_id is required")
	}

	if body.Marks < 0 || body.Marks > 20 {
		return ReturnBadRequest(c, "marks must be between 0 and 20")
	}

	if body.Date == "" {
		return ReturnBadRequest(c, "date is required")
	}
	parsedDate, err := time.Parse("2006-01-02", body.Date)
	if err != nil {
		return ReturnBadRequest(c, "invalid date format, use yyyy-mm-dd")
	}

	var registration db.Registration
	if err := db.DB().Preload("StudentClass").First(&registration, registrationId).Error; err != nil {
		return ReturnNotFound(c, "registration not found")
	}

	var category db.EvaluationCategory
	if err := db.DB().First(&category, body.EvaluationCategoryID).Error; err != nil {
		return ReturnNotFound(c, "evaluation category not found")
	}

	var existingCount int64
	db.DB().Model(&db.Grade{}).Where(
		"category_id = ? AND registration_id = ?",
		body.EvaluationCategoryID, registrationId,
	).Count(&existingCount)

	if int(existingCount) >= category.Cardinality {
		return ReturnBadRequest(c, fmt.Sprintf("cannot create grade: all %d expected evaluations for this category are already registered", category.Cardinality))
	}

	userEmail, err := GetUserEmailFromToken(c)
	if err != nil {
		return ReturnUnauthorized(c, err.Error())
	}

	calculator := services.NewGradingCalculator()
	percentage := calculator.CalculatePercentage(body.Marks, category.MaxScore)

	now := time.Now()
	evalNumber := int(existingCount) + 1
	name := fmt.Sprintf("%s %d", category.Name, evalNumber)

	grade := db.Grade{
		CategoryID:     body.EvaluationCategoryID,
		StudentClassID: registration.StudentClassID,
		RegistrationID: registrationId,
		Name:           name,
		Date:           &parsedDate,
		MaxScore:       category.MaxScore,
		Score:          &body.Marks,
		Percentage:     &percentage,
		GradedAt:       &now,
		CreatedAt:      now,
		CreatedBy:      &userEmail,
	}

	if err := db.DB().Create(&grade).Error; err != nil {
		return ReturnInternalError(c, "failed to create grade")
	}

	db.DB().Preload("Category").Preload("StudentClass").First(&grade, grade.ID)

	return c.Status(201).JSON(grade)
}

func UpdateGrade(c *fiber.Ctx) error {
	parsedRegID, err := strconv.ParseUint(c.Params("registrationId"), 10, 64)
	if err != nil {
		return ReturnBadRequest(c, "invalid registrationId")
	}
	registrationId := uint(parsedRegID)

	parsedGradeID, err := strconv.ParseUint(c.Params("gradeId"), 10, 64)
	if err != nil {
		return ReturnBadRequest(c, "invalid gradeId")
	}
	gradeId := uint(parsedGradeID)

	var body struct {
		Marks *float64 `json:"marks"`
		Date  *string  `json:"date"`
	}

	if err := c.BodyParser(&body); err != nil {
		return ReturnBadRequest(c, "invalid request body")
	}

	if body.Marks == nil && body.Date == nil {
		return ReturnBadRequest(c, "at least one of marks or date is required")
	}

	if body.Marks != nil && (*body.Marks < 0 || *body.Marks > 20) {
		return ReturnBadRequest(c, "marks must be between 0 and 20")
	}

	var parsedDate *time.Time
	if body.Date != nil {
		d, err := time.Parse("2006-01-02", *body.Date)
		if err != nil {
			return ReturnBadRequest(c, "invalid date format, use yyyy-mm-dd")
		}
		parsedDate = &d
	}

	var grade db.Grade
	if err := db.DB().Preload("Category").First(&grade, gradeId).Error; err != nil {
		return ReturnNotFound(c, "grade not found")
	}

	if grade.RegistrationID != registrationId {
		return ReturnNotFound(c, "grade not found for this registration")
	}

	userEmail, err := GetUserEmailFromToken(c)
	if err != nil {
		return ReturnUnauthorized(c, err.Error())
	}

	if body.Marks != nil {
		calculator := services.NewGradingCalculator()
		grade.Score = body.Marks
		if !grade.IsExcused {
			percentage := calculator.CalculatePercentage(*body.Marks, grade.MaxScore)
			grade.Percentage = &percentage
		}
		now := time.Now()
		grade.GradedAt = &now
	}

	if parsedDate != nil {
		grade.Date = parsedDate
	}

	grade.UpdatedBy = &userEmail

	if err := db.DB().Save(&grade).Error; err != nil {
		return ReturnInternalError(c, "failed to update grade")
	}

	db.DB().Preload("Category").Preload("StudentClass").First(&grade, grade.ID)

	return c.JSON(grade)
}

func BatchCreateGrades(c *fiber.Ctx) error {
	var request struct {
		CategoryID     uint                     `json:"category_id"`
		StudentClassID uint                     `json:"student_class_id"`
		Name           string                   `json:"name"`
		Description    *string                  `json:"description"`
		Date           *time.Time               `json:"date"`
		DueDate        *time.Time               `json:"due_date"`
		MaxScore       float64                  `json:"max_score"`
		Grades         []map[string]interface{} `json:"grades"`
	}

	if err := c.BodyParser(&request); err != nil {
		return ReturnBadRequest(c, "invalid request body")
	}

	if request.CategoryID == 0 {
		return ReturnBadRequest(c, "category_id is required")
	}

	if request.StudentClassID == 0 {
		return ReturnBadRequest(c, "student_class_id is required")
	}

	if request.Name == "" {
		return ReturnBadRequest(c, "name is required")
	}

	if request.MaxScore <= 0 {
		return ReturnBadRequest(c, "max_score must be greater than 0")
	}

	var studentClass db.StudentClass
	if err := db.DB().First(&studentClass, request.StudentClassID).Error; err != nil {
		return ReturnNotFound(c, "student class not found")
	}

	var category db.EvaluationCategory
	if err := db.DB().First(&category, request.CategoryID).Error; err != nil {
		return ReturnNotFound(c, "evaluation category not found")
	}

	if category.CourseID != studentClass.CourseID {
		return ReturnBadRequest(c, "evaluation category course does not match student class course")
	}

	registrationGradeCounts := make(map[uint]int)
	for _, gradeData := range request.Grades {
		registrationID := uint(gradeData["registration_id"].(float64))
		registrationGradeCounts[registrationID]++
	}

	for regID, newCount := range registrationGradeCounts {
		var existingCount int64
		db.DB().Model(&db.Grade{}).Where("category_id = ? AND student_class_id = ? AND registration_id = ?",
			request.CategoryID, request.StudentClassID, regID).Count(&existingCount)

		if int(existingCount)+newCount > category.Cardinality {
			return ReturnBadRequest(c, fmt.Sprintf("cannot create grades: would exceed cardinality limit (%d) for registration %d", category.Cardinality, regID))
		}
	}

	calculator := services.NewGradingCalculator()
	var createdGrades []db.Grade
	registrationIDsToUpdate := make(map[uint]bool)

	for _, gradeData := range request.Grades {
		registrationID := uint(gradeData["registration_id"].(float64))
		score := gradeData["score"].(float64)

		percentage := calculator.CalculatePercentage(score, request.MaxScore)

		grade := db.Grade{
			CategoryID:     request.CategoryID,
			StudentClassID: request.StudentClassID,
			RegistrationID: registrationID,
			Name:           request.Name,
			Description:    request.Description,
			Date:           request.Date,
			DueDate:        request.DueDate,
			MaxScore:       request.MaxScore,
			Score:          &score,
			Percentage:     &percentage,
		}

		if comments, ok := gradeData["comments"].(string); ok {
			grade.Comments = &comments
		}
		if gradedBy, ok := gradeData["graded_by"].(string); ok {
			grade.GradedBy = &gradedBy
		}

		if err := db.DB().Create(&grade).Error; err == nil {
			createdGrades = append(createdGrades, grade)
			registrationIDsToUpdate[registrationID] = true
		}
	}

	return c.Status(201).JSON(fiber.Map{
		"created": len(createdGrades),
		"grades":  createdGrades,
	})
}

func GetGradeStatistics(c *fiber.Ctx) error {
	categoryID, err := ParseOptionalUintQueryParam(c, "category_id")
	if err != nil {
		return ReturnBadRequest(c, err.Error())
	}

	studentClassID, err := ParseOptionalUintQueryParam(c, "student_class_id")
	if err != nil {
		return ReturnBadRequest(c, err.Error())
	}

	registrationID, err := ParseOptionalUintQueryParam(c, "registration_id")
	if err != nil {
		return ReturnBadRequest(c, err.Error())
	}

	if studentClassID == nil {
		return ReturnBadRequest(c, "student_class_id is required")
	}

	if registrationID == nil {
		return ReturnBadRequest(c, "registration_id is required")
	}

	calculator := services.NewGradingCalculator()

	if categoryID != nil {
		stats, err := calculator.GetCategoryStatistics(*categoryID, *studentClassID, *registrationID)
		if err != nil {
			return ReturnInternalError(c, "failed to fetch statistics")
		}
		return c.JSON(stats)
	}

	var sc db.StudentClass
	if err := db.DB().First(&sc, *studentClassID).Error; err != nil {
		return ReturnInternalError(c, "student class not found")
	}

	allStats, err := calculator.GetAllCategoryStatistics(*studentClassID, *registrationID, sc.CourseID, sc.LevelId)
	if err != nil {
		return ReturnInternalError(c, "failed to fetch all statistics")
	}

	return c.JSON(allStats)
}
