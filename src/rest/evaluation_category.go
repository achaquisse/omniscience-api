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
