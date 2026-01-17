package rest

import (
	"skulla-api/db"

	"github.com/gofiber/fiber/v2"
)

func ListStudentClass(c *fiber.Ctx) error {
	startDate, err := ParseDateQueryParam(c, "startDate")
	if err != nil {
		return ReturnBadRequest(c, err.Error())
	}

	endDate, err := ParseDateQueryParam(c, "endDate")
	if err != nil {
		return ReturnBadRequest(c, err.Error())
	}

	userEmail, err := GetUserEmailFromToken(c)
	if err != nil {
		return ReturnBadRequest(c, err.Error())
	}

	var courses = db.ListCoursesByTeacherEmail(userEmail)
	var courseIds []uint
	for _, course := range courses {
		courseIds = append(courseIds, course.ID)
	}

	studentClass := db.ListStudentClasses(courseIds, startDate, endDate)
	return c.JSON(studentClass)
}
