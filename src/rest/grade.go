package rest

import (
	"skulla-api/db"
	"skulla-api/services"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func ListGrades(c *fiber.Ctx) error {
	evaluationItemIDStr := c.Query("evaluation_item_id")
	registrationIDStr := c.Query("registration_id")
	studentClassIDStr := c.Query("student_class_id")

	query := db.DB().Preload("EvaluationItem").Preload("EvaluationItem.Category")

	if evaluationItemIDStr != "" {
		evaluationItemID, err := strconv.ParseUint(evaluationItemIDStr, 10, 64)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid evaluation_item_id"})
		}
		query = query.Where("evaluation_item_id = ?", evaluationItemID)
	}

	if registrationIDStr != "" {
		registrationID, err := strconv.ParseUint(registrationIDStr, 10, 64)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid registration_id"})
		}
		query = query.Where("registration_id = ?", registrationID)
	}

	if studentClassIDStr != "" {
		studentClassID, err := strconv.ParseUint(studentClassIDStr, 10, 64)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid student_class_id"})
		}
		query = query.Joins("JOIN EvaluationItem ON EvaluationItem.id = Grade.evaluation_item_id").
			Where("EvaluationItem.student_class_id = ?", studentClassID)
	}

	var grades []db.Grade
	err := query.Find(&grades).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch grades"})
	}

	return c.JSON(grades)
}

func GetGrade(c *fiber.Ctx) error {
	id := c.Params("id")
	var grade db.Grade

	err := db.DB().Preload("EvaluationItem").Preload("EvaluationItem.Category").
		First(&grade, id).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "grade not found"})
	}

	return c.JSON(grade)
}

func CreateGrade(c *fiber.Ctx) error {
	var grade db.Grade

	if err := c.BodyParser(&grade); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	if grade.EvaluationItemID == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "evaluation_item_id is required"})
	}

	if grade.RegistrationID == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "registration_id is required"})
	}

	var evaluationItem db.EvaluationItem
	err := db.DB().First(&evaluationItem, grade.EvaluationItemID).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "evaluation item not found"})
	}

	if grade.Score != nil && !grade.IsExcused {
		calculator := services.NewGradingCalculator()
		percentage := calculator.CalculatePercentage(*grade.Score, evaluationItem.MaxScore)
		grade.Percentage = &percentage
	}

	err = db.DB().Create(&grade).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to create grade"})
	}

	db.DB().Preload("EvaluationItem").Preload("EvaluationItem.Category").First(&grade, grade.ID)

	return c.Status(201).JSON(grade)
}

func UpdateGrade(c *fiber.Ctx) error {
	id := c.Params("id")
	var grade db.Grade

	err := db.DB().Preload("EvaluationItem").First(&grade, id).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "grade not found"})
	}

	oldScore := grade.Score
	oldStatus := grade.Status

	var updates db.Grade
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	if updates.Score != nil {
		grade.Score = updates.Score
		if !grade.IsExcused {
			calculator := services.NewGradingCalculator()
			percentage := calculator.CalculatePercentage(*updates.Score, grade.EvaluationItem.MaxScore)
			grade.Percentage = &percentage
		}
	}

	grade.IsExcused = updates.IsExcused
	grade.IsLate = updates.IsLate
	grade.LatePenalty = updates.LatePenalty
	if updates.Comments != nil {
		grade.Comments = updates.Comments
	}
	if updates.GradedBy != nil {
		grade.GradedBy = updates.GradedBy
	}
	if updates.Status != "" {
		grade.Status = updates.Status
		if updates.Status == "PUBLISHED" && grade.GradedAt == nil {
			now := time.Now()
			grade.GradedAt = &now
		}
	}
	if updates.UpdatedBy != nil {
		grade.UpdatedBy = updates.UpdatedBy
	}

	err = db.DB().Save(&grade).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to update grade"})
	}

	if oldScore != updates.Score || oldStatus != updates.Status {
		history := db.GradeHistory{
			GradeID:   grade.ID,
			OldScore:  oldScore,
			NewScore:  updates.Score,
			OldStatus: &oldStatus,
			NewStatus: &updates.Status,
			ChangedBy: *updates.UpdatedBy,
		}
		db.DB().Create(&history)
	}

	db.DB().Preload("EvaluationItem").Preload("EvaluationItem.Category").First(&grade, grade.ID)

	return c.JSON(grade)
}

func BatchCreateGrades(c *fiber.Ctx) error {
	var request struct {
		EvaluationItemID uint                     `json:"evaluation_item_id"`
		Grades           []map[string]interface{} `json:"grades"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	if request.EvaluationItemID == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "evaluation_item_id is required"})
	}

	var evaluationItem db.EvaluationItem
	err := db.DB().First(&evaluationItem, request.EvaluationItemID).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "evaluation item not found"})
	}

	calculator := services.NewGradingCalculator()
	var createdGrades []db.Grade

	for _, gradeData := range request.Grades {
		registrationID := uint(gradeData["registration_id"].(float64))
		score := gradeData["score"].(float64)

		percentage := calculator.CalculatePercentage(score, evaluationItem.MaxScore)

		grade := db.Grade{
			EvaluationItemID: request.EvaluationItemID,
			RegistrationID:   registrationID,
			Score:            &score,
			Percentage:       &percentage,
			Status:           "DRAFT",
		}

		if comments, ok := gradeData["comments"].(string); ok {
			grade.Comments = &comments
		}
		if gradedBy, ok := gradeData["graded_by"].(string); ok {
			grade.GradedBy = &gradedBy
		}

		err = db.DB().Create(&grade).Error
		if err == nil {
			createdGrades = append(createdGrades, grade)
		}
	}

	return c.Status(201).JSON(fiber.Map{
		"created": len(createdGrades),
		"grades":  createdGrades,
	})
}

func PublishGrades(c *fiber.Ctx) error {
	evaluationItemIDStr := c.Params("evaluation_item_id")
	evaluationItemID, err := strconv.ParseUint(evaluationItemIDStr, 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid evaluation_item_id"})
	}

	now := time.Now()
	result := db.DB().Model(&db.Grade{}).
		Where("evaluation_item_id = ? AND status = ?", evaluationItemID, "DRAFT").
		Updates(map[string]interface{}{
			"status":    "PUBLISHED",
			"graded_at": now,
		})

	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to publish grades"})
	}

	return c.JSON(fiber.Map{
		"message":   "grades published successfully",
		"published": result.RowsAffected,
	})
}
