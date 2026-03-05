package rest

import (
	"skulla-api/db"
	"skulla-api/services"
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

	levelIDStr := c.Query("level_id")

	query := db.DB().Where("course_id = ? AND is_active = ?", courseID, true)

	if levelIDStr != "" {
		levelID, err := strconv.ParseUint(levelIDStr, 10, 64)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid level_id"})
		}
		query = query.Where("level_id = ?", levelID)
	} else {
		query = query.Where("level_id IS NULL")
	}

	var formula db.GradingFormula
	err = query.Preload("Course").First(&formula).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "grading formula not found"})
	}

	return c.JSON(formula)
}

func GetFormulaResults(c *fiber.Ctx) error {
	registrationID, err := ParseOptionalUintQueryParam(c, "registration_id")
	if err != nil {
		return ReturnBadRequest(c, err.Error())
	}
	if registrationID == nil {
		return ReturnBadRequest(c, "registration_id is required")
	}

	var registration db.Registration
	err = db.DB().First(&registration, *registrationID).Error
	if err != nil {
		return ReturnInternalError(c, "registration not found")
	}

	studentClassID := registration.StudentClassID

	courseID, err := db.GetStudentClassCourseID(studentClassID)
	if err != nil {
		return ReturnInternalError(c, "failed to get course ID")
	}

	calculator := services.NewGradingCalculator()
	result, err := calculator.CalculateFormulaDetailed(*registrationID, studentClassID, courseID)
	if err != nil {
		return ReturnInternalError(c, err.Error())
	}

	return c.JSON(result)
}
