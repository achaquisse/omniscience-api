package rest

import (
	"skulla-api/db"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetGradingFormula(c *fiber.Ctx) error {
	courseIDStr := c.Query("course_id")
	if courseIDStr == "" {
		return c.Status(400).JSON(fiber.Map{"error": "course_id is required"})
	}

	courseID, err := strconv.ParseUint(courseIDStr, 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid course_id"})
	}

	var formula db.GradingFormula
	err = db.DB().Where("course_id = ? AND is_active = ?", courseID, true).
		Preload("Course").First(&formula).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "grading formula not found"})
	}

	return c.JSON(formula)
}

func CreateGradingFormula(c *fiber.Ctx) error {
	var formula db.GradingFormula

	if err := c.BodyParser(&formula); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	if formula.CourseID == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "course_id is required"})
	}

	if formula.FormulaType == "" {
		return c.Status(400).JSON(fiber.Map{"error": "formula_type is required"})
	}

	err := db.DB().Create(&formula).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to create formula"})
	}

	db.DB().Preload("Course").First(&formula, formula.ID)

	return c.Status(201).JSON(formula)
}

func UpdateGradingFormula(c *fiber.Ctx) error {
	id := c.Params("id")
	var formula db.GradingFormula

	err := db.DB().First(&formula, id).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "formula not found"})
	}

	var updates db.GradingFormula
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	if updates.FormulaType != "" {
		formula.FormulaType = updates.FormulaType
	}
	if updates.FormulaConfig != nil {
		formula.FormulaConfig = updates.FormulaConfig
	}
	if updates.PassingPercentage != nil {
		formula.PassingPercentage = updates.PassingPercentage
	}
	if updates.GradingScale != nil {
		formula.GradingScale = updates.GradingScale
	}
	formula.IsActive = updates.IsActive
	if updates.UpdatedBy != nil {
		formula.UpdatedBy = updates.UpdatedBy
	}

	err = db.DB().Save(&formula).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to update formula"})
	}

	db.DB().Preload("Course").First(&formula, formula.ID)

	return c.JSON(formula)
}
