package rest

import (
	"fmt"
	"skulla-api/db"
	"skulla-api/services"
	"time"

	"github.com/gofiber/fiber/v2"
)

func ListGrades(c *fiber.Ctx) error {
	query := db.DB().Preload("Category").Preload("StudentClass")

	if categoryID, err := ParseOptionalUintQueryParam(c, "category_id"); err != nil {
		return ReturnBadRequest(c, err.Error())
	} else if categoryID != nil {
		query = query.Where("category_id = ?", *categoryID)
	}

	if registrationID, err := ParseOptionalUintQueryParam(c, "registration_id"); err != nil {
		return ReturnBadRequest(c, err.Error())
	} else if registrationID != nil {
		query = query.Where("registration_id = ?", *registrationID)
	}

	if studentClassID, err := ParseOptionalUintQueryParam(c, "student_class_id"); err != nil {
		return ReturnBadRequest(c, err.Error())
	} else if studentClassID != nil {
		query = query.Where("student_class_id = ?", *studentClassID)
	}

	var grades []db.Grade
	if err := query.Find(&grades).Error; err != nil {
		return ReturnInternalError(c, "failed to fetch grades")
	}

	return c.JSON(grades)
}

func GetGrade(c *fiber.Ctx) error {
	id := c.Params("id")
	var grade db.Grade

	if err := db.DB().Preload("Category").Preload("StudentClass").First(&grade, id).Error; err != nil {
		return ReturnNotFound(c, "grade not found")
	}

	return c.JSON(grade)
}

func CreateGrade(c *fiber.Ctx) error {
	var grade db.Grade

	if err := c.BodyParser(&grade); err != nil {
		return ReturnBadRequest(c, "invalid request body")
	}

	if grade.CategoryID == 0 {
		return ReturnBadRequest(c, "category_id is required")
	}

	if grade.RegistrationID == 0 {
		return ReturnBadRequest(c, "registration_id is required")
	}

	var registration db.Registration
	if err := db.DB().Preload("StudentClass").First(&registration, grade.RegistrationID).Error; err != nil {
		return ReturnNotFound(c, "registration not found")
	}

	var category db.EvaluationCategory
	if err := db.DB().First(&category, grade.CategoryID).Error; err != nil {
		return ReturnNotFound(c, "evaluation category not found")
	}

	if category.CourseID != registration.StudentClass.CourseID {
		return ReturnBadRequest(c, "evaluation category course does not match student class course")
	}

	var existingGradesCount int64
	db.DB().Model(&db.Grade{}).Where("category_id = ? AND student_class_id = ? AND registration_id = ?",
		grade.CategoryID, registration.StudentClassID, grade.RegistrationID).Count(&existingGradesCount)

	if int(existingGradesCount) >= category.Cardinality {
		return ReturnBadRequest(c, fmt.Sprintf("cannot create more grades: category cardinality limit (%d) reached", category.Cardinality))
	}

	if grade.Score != nil {
		if *grade.Score < 0 {
			return ReturnBadRequest(c, "score cannot be negative")
		}
		if *grade.Score > category.MaxScore {
			return ReturnBadRequest(c, fmt.Sprintf("score cannot exceed max_score: %g", category.MaxScore))
		}
	}

	grade.StudentClassID = registration.StudentClassID
	grade.Name = category.Name
	grade.MaxScore = category.MaxScore

	calculator := services.NewGradingCalculator()
	if grade.Score != nil && !grade.IsExcused {
		percentage := calculator.CalculatePercentage(*grade.Score, grade.MaxScore)
		grade.Percentage = &percentage
	}

	if err := db.DB().Create(&grade).Error; err != nil {
		return ReturnInternalError(c, "failed to create grade")
	}

	db.DB().Preload("Category").Preload("StudentClass").First(&grade, grade.ID)

	return c.Status(201).JSON(grade)
}

func UpdateGrade(c *fiber.Ctx) error {
	id := c.Params("id")
	var grade db.Grade

	if err := db.DB().Preload("Category").First(&grade, id).Error; err != nil {
		return ReturnNotFound(c, "grade not found")
	}

	oldScore := grade.Score

	var updates struct {
		Score     *float64 `json:"score"`
		UpdatedBy *string  `json:"updated_by,omitempty"`
	}

	if err := c.BodyParser(&updates); err != nil {
		return ReturnBadRequest(c, "invalid request body")
	}

	if updates.Score == nil {
		return ReturnBadRequest(c, "score is required")
	}

	if *updates.Score < 0 {
		return ReturnBadRequest(c, "score cannot be negative")
	}

	if *updates.Score > grade.Category.MaxScore {
		return ReturnBadRequest(c, fmt.Sprintf("score cannot exceed max_score: %g", grade.Category.MaxScore))
	}

	calculator := services.NewGradingCalculator()
	grade.Score = updates.Score

	if !grade.IsExcused {
		percentage := calculator.CalculatePercentage(*updates.Score, grade.MaxScore)
		grade.Percentage = &percentage
	}

	if grade.GradedAt == nil {
		now := time.Now()
		grade.GradedAt = &now
	}

	if updates.UpdatedBy != nil {
		grade.UpdatedBy = updates.UpdatedBy
	}

	if err := db.DB().Save(&grade).Error; err != nil {
		return ReturnInternalError(c, "failed to update grade")
	}

	if oldScore != updates.Score {
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

	courseID, err := db.GetStudentClassCourseID(*studentClassID)
	if err != nil {
		return ReturnInternalError(c, "failed to get course ID")
	}

	allStats, err := calculator.GetAllCategoryStatistics(*studentClassID, *registrationID, courseID)
	if err != nil {
		return ReturnInternalError(c, "failed to fetch all statistics")
	}

	return c.JSON(allStats)
}
