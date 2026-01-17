package rest

import (
	"skulla-api/db"

	"github.com/gofiber/fiber/v2"
)

func ListRegistrations(c *fiber.Ctx) error {
	studentClassId, err := ParseUintQueryParam(c, "studentClassId", true)
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

	if db.IsTeacherEmailBelongToCourse(userEmail, int(courseID)) {
		registrations := db.ListRegistrations(int(studentClassId))
		return c.JSON(registrations)
	}

	return ReturnUnauthorized(c, "User does not have permission to access course")
}
