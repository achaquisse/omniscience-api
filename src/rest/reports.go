package rest

import (
	"skulla-api/db"
	"skulla-api/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetStudentReport(c *fiber.Ctx) error {
	registrationIDStr := c.Params("registration_id")
	registrationID, err := strconv.ParseUint(registrationIDStr, 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid registration_id"})
	}

	var registration db.Registration
	err = db.DB().First(&registration, registrationID).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "registration not found"})
	}

	var grades []db.Grade
	err = db.DB().
		Preload("EvaluationItem").
		Preload("EvaluationItem.Category").
		Joins("JOIN EvaluationItem ON EvaluationItem.id = Grade.evaluation_item_id").
		Where("Grade.registration_id = ? AND Grade.status = ?", registrationID, "PUBLISHED").
		Find(&grades).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch grades"})
	}

	categoryScores := make(map[string]map[string]interface{})
	for _, grade := range grades {
		categoryName := grade.EvaluationItem.Category.Name
		if _, exists := categoryScores[categoryName]; !exists {
			categoryScores[categoryName] = map[string]interface{}{
				"category_id":   grade.EvaluationItem.Category.ID,
				"category_name": categoryName,
				"weight":        grade.EvaluationItem.Category.Weight,
				"grades":        []db.Grade{},
			}
		}
		gradesList := categoryScores[categoryName]["grades"].([]db.Grade)
		categoryScores[categoryName]["grades"] = append(gradesList, grade)
	}

	var finalGrade db.FinalGrade
	err = db.DB().
		Where("registration_id = ? AND student_class_id = ?", registrationID, registration.StudentClassID).
		First(&finalGrade).Error

	report := fiber.Map{
		"registration_id":  registrationID,
		"student_class_id": registration.StudentClassID,
		"category_scores":  categoryScores,
		"total_grades":     len(grades),
	}

	if err == nil {
		report["final_grade"] = finalGrade
	}

	return c.JSON(report)
}

func GetClassReport(c *fiber.Ctx) error {
	studentClassIDStr := c.Params("student_class_id")
	studentClassID, err := strconv.ParseUint(studentClassIDStr, 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid student_class_id"})
	}

	var studentClass db.StudentClass
	err = db.DB().Preload("Course").First(&studentClass, studentClassID).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "student class not found"})
	}

	type PerformanceStats struct {
		EvaluationItemID  uint    `json:"evaluation_item_id"`
		EvaluationName    string  `json:"evaluation_name"`
		CategoryName      string  `json:"category_name"`
		TotalGraded       int64   `json:"total_graded"`
		AveragePercentage float64 `json:"average_percentage"`
		MinPercentage     float64 `json:"min_percentage"`
		MaxPercentage     float64 `json:"max_percentage"`
		PassingCount      int64   `json:"passing_count"`
		FailingCount      int64   `json:"failing_count"`
	}

	var stats []PerformanceStats
	err = db.DB().Table("Grade").
		Select(`
			EvaluationItem.id as evaluation_item_id,
			EvaluationItem.name as evaluation_name,
			EvaluationCategory.name as category_name,
			COUNT(DISTINCT Grade.id) as total_graded,
			AVG(Grade.percentage) as average_percentage,
			MIN(Grade.percentage) as min_percentage,
			MAX(Grade.percentage) as max_percentage,
			SUM(CASE WHEN Grade.percentage >= 60 THEN 1 ELSE 0 END) as passing_count,
			SUM(CASE WHEN Grade.percentage < 60 THEN 1 ELSE 0 END) as failing_count
		`).
		Joins("INNER JOIN EvaluationItem ON EvaluationItem.id = Grade.evaluation_item_id").
		Joins("INNER JOIN EvaluationCategory ON EvaluationCategory.id = EvaluationItem.category_id").
		Where("EvaluationItem.student_class_id = ? AND Grade.status = ? AND Grade.is_excused = ?",
			studentClassID, "PUBLISHED", false).
		Group("EvaluationItem.id, EvaluationCategory.id").
		Scan(&stats).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch class statistics"})
	}

	var finalGrades []db.FinalGrade
	err = db.DB().Where("student_class_id = ? AND status = ?", studentClassID, "PUBLISHED").
		Find(&finalGrades).Error

	var avgFinalGrade float64
	var passingStudents int64
	var totalStudents int64 = int64(len(finalGrades))

	for _, fg := range finalGrades {
		if fg.CalculatedPercentage != nil {
			avgFinalGrade += *fg.CalculatedPercentage
		}
		if fg.IsPassing != nil && *fg.IsPassing {
			passingStudents++
		}
	}

	if totalStudents > 0 {
		avgFinalGrade = avgFinalGrade / float64(totalStudents)
	}

	return c.JSON(fiber.Map{
		"student_class_id":    studentClassID,
		"class_name":          studentClass.Name,
		"course_name":         studentClass.Course.Name,
		"evaluation_stats":    stats,
		"total_students":      totalStudents,
		"passing_students":    passingStudents,
		"average_final_grade": avgFinalGrade,
		"final_grades":        finalGrades,
	})
}

func CalculateFinalGrades(c *fiber.Ctx) error {
	studentClassIDStr := c.Params("student_class_id")
	studentClassID, err := strconv.ParseUint(studentClassIDStr, 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid student_class_id"})
	}

	calculator := services.NewGradingCalculator()
	err = calculator.RecalculateAllFinalGrades(uint(studentClassID))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to calculate final grades: " + err.Error()})
	}

	return c.JSON(fiber.Map{
		"message":          "final grades calculated successfully",
		"student_class_id": studentClassID,
	})
}

func GetFinalGrade(c *fiber.Ctx) error {
	registrationIDStr := c.Params("registration_id")
	registrationID, err := strconv.ParseUint(registrationIDStr, 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid registration_id"})
	}

	var finalGrade db.FinalGrade
	err = db.DB().
		Preload("StudentClass").
		Preload("StudentClass.Course").
		Where("registration_id = ?", registrationID).
		First(&finalGrade).Error

	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "final grade not found"})
	}

	return c.JSON(finalGrade)
}
