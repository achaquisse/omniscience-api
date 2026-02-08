package rest

import (
	"skulla-api/db"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func ListEvaluationCategories(c *fiber.Ctx) error {
	courseIDStr := c.Query("course_id")
	if courseIDStr == "" {
		return c.Status(400).JSON(fiber.Map{"error": "course_id is required"})
	}

	courseID, err := strconv.ParseUint(courseIDStr, 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid course_id"})
	}

	var categories []db.EvaluationCategory
	query := db.DB().Where("course_id = ?", courseID)

	includeInactive := c.Query("include_inactive")
	if includeInactive != "true" {
		query = query.Where("is_active = ?", true)
	}

	err = query.Order("display_order ASC").Preload("Course").Find(&categories).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch categories"})
	}

	return c.JSON(categories)
}

func GetEvaluationCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	var category db.EvaluationCategory

	err := db.DB().Preload("Course").First(&category, id).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "category not found"})
	}

	return c.JSON(category)
}

func CreateEvaluationCategory(c *fiber.Ctx) error {
	var category db.EvaluationCategory

	if err := c.BodyParser(&category); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	if category.CourseID == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "course_id is required"})
	}

	if category.Name == "" {
		return c.Status(400).JSON(fiber.Map{"error": "name is required"})
	}

	if category.Weight < 0 || category.Weight > 100 {
		return c.Status(400).JSON(fiber.Map{"error": "weight must be between 0 and 100"})
	}

	err := db.DB().Create(&category).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to create category"})
	}

	db.DB().Preload("Course").First(&category, category.ID)

	return c.Status(201).JSON(category)
}

func UpdateEvaluationCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	var category db.EvaluationCategory

	err := db.DB().First(&category, id).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "category not found"})
	}

	var updates db.EvaluationCategory
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	if updates.Name != "" {
		category.Name = updates.Name
	}
	if updates.Description != nil {
		category.Description = updates.Description
	}
	if updates.Weight >= 0 && updates.Weight <= 100 {
		category.Weight = updates.Weight
	}
	if updates.MaxScore > 0 {
		category.MaxScore = updates.MaxScore
	}
	if updates.DropLowest >= 0 {
		category.DropLowest = updates.DropLowest
	}
	category.IsExtraCredit = updates.IsExtraCredit
	category.DisplayOrder = updates.DisplayOrder
	category.IsActive = updates.IsActive

	if updates.UpdatedBy != nil {
		category.UpdatedBy = updates.UpdatedBy
	}

	err = db.DB().Save(&category).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to update category"})
	}

	db.DB().Preload("Course").First(&category, category.ID)

	return c.JSON(category)
}

func DeleteEvaluationCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	var category db.EvaluationCategory

	err := db.DB().First(&category, id).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "category not found"})
	}

	err = db.DB().Delete(&category).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to delete category"})
	}

	return c.Status(204).Send(nil)
}
