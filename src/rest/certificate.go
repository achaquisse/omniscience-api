package rest

import (
	"skulla-api/db"
	"skulla-api/services"

	"github.com/gofiber/fiber/v2"
)

func GenerateCertificate(c *fiber.Ctx) error {
	template := c.Query("template")
	studentName := c.Query("studentName")
	certDescription := c.Query("certDescription")

	if template != "portuguese" && template != "english" {
		return ReturnBadRequest(c, "template must be 'portuguese' or 'english'")
	}

	if len(studentName) < 5 {
		return ReturnBadRequest(c, "studentName must be at least 5 characters")
	}

	if len(certDescription) < 10 {
		return ReturnBadRequest(c, "certDescription must be at least 10 characters")
	}

	pdfBytes, err := services.GenerateCertificate(template, studentName, certDescription)
	if err != nil {
		return ReturnBadRequest(c, err.Error())
	}

	userEmail, err := GetUserEmailFromToken(c)
	if err != nil {
		return ReturnInternalError(c, err.Error())
	}

	if err := db.CreateCertificateLog(template, studentName, certDescription, userEmail); err != nil {
		return ReturnInternalError(c, err.Error())
	}

	c.Set("Content-Type", "application/pdf")
	return c.Send(pdfBytes)
}
