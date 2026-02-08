package rest

import (
	"skulla-api/db"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func ListEvaluationItems(c *fiber.Ctx) error {
	categoryIDStr := c.Query("category_id")
	studentClassIDStr := c.Query("student_class_id")

	query := db.DB().Preload("Category").Preload("StudentClass")

	if categoryIDStr != "" {
		categoryID, err := strconv.ParseUint(categoryIDStr, 10, 64)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid category_id"})
		}
		query = query.Where("category_id = ?", categoryID)
	}

	if studentClassIDStr != "" {
		studentClassID, err := strconv.ParseUint(studentClassIDStr, 10, 64)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid student_class_id"})
		}
		query = query.Where("student_class_id = ?", studentClassID)
	}

	var items []db.EvaluationItem
	err := query.Order("date DESC").Find(&items).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch evaluation items"})
	}

	return c.JSON(items)
}

func GetEvaluationItem(c *fiber.Ctx) error {
	id := c.Params("id")
	var item db.EvaluationItem

	err := db.DB().Preload("Category").Preload("StudentClass").First(&item, id).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "evaluation item not found"})
	}

	return c.JSON(item)
}

func CreateEvaluationItem(c *fiber.Ctx) error {
	var item db.EvaluationItem

	if err := c.BodyParser(&item); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	if item.CategoryID == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "category_id is required"})
	}

	if item.StudentClassID == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "student_class_id is required"})
	}

	if item.Name == "" {
		return c.Status(400).JSON(fiber.Map{"error": "name is required"})
	}

	if item.MaxScore <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "max_score must be greater than 0"})
	}

	err := db.DB().Create(&item).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to create evaluation item"})
	}

	db.DB().Preload("Category").Preload("StudentClass").First(&item, item.ID)

	return c.Status(201).JSON(item)
}

func UpdateEvaluationItem(c *fiber.Ctx) error {
	id := c.Params("id")
	var item db.EvaluationItem

	err := db.DB().First(&item, id).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "evaluation item not found"})
	}

	var updates db.EvaluationItem
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	if updates.Name != "" {
		item.Name = updates.Name
	}
	if updates.Description != nil {
		item.Description = updates.Description
	}
	if updates.Date != nil {
		item.Date = updates.Date
	}
	if updates.DueDate != nil {
		item.DueDate = updates.DueDate
	}
	if updates.MaxScore > 0 {
		item.MaxScore = updates.MaxScore
	}
	if updates.WeightOverride != nil {
		item.WeightOverride = updates.WeightOverride
	}
	if updates.Status != "" {
		item.Status = updates.Status
	}
	if updates.UpdatedBy != nil {
		item.UpdatedBy = updates.UpdatedBy
	}

	err = db.DB().Save(&item).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to update evaluation item"})
	}

	db.DB().Preload("Category").Preload("StudentClass").First(&item, item.ID)

	return c.JSON(item)
}

func DeleteEvaluationItem(c *fiber.Ctx) error {
	id := c.Params("id")
	var item db.EvaluationItem

	err := db.DB().First(&item, id).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "evaluation item not found"})
	}

	err = db.DB().Delete(&item).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to delete evaluation item"})
	}

	return c.Status(204).Send(nil)
}
