package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func Init(app *fiber.App) {
	SetupSwagger(app)

	app.Get("/student-classes", AuthMiddleware, ListStudentClass)
	app.Get("/registrations", AuthMiddleware, ListRegistrations)

	app.Post("/attendance", AuthMiddleware, RecordAttendance)
	app.Post("/attendance/bulk", AuthMiddleware, RecordBulkAttendance)
	app.Get("/attendance/report", AuthMiddleware, GetStudentAttendanceReport)
	app.Get("/attendance/class-report", AuthMiddleware, GetClassAttendanceReport)
	app.Post("/attendance/trigger-reports", AuthMiddleware, TriggerAttendanceReports)

	app.Get("/grades/categories", AuthMiddleware, ListEvaluationCategories)
	app.Get("/grades/formulas", AuthMiddleware, GetGradingFormula)
	app.Get("/grades/formulas/results", AuthMiddleware, GetFormulaResults)
	app.Get("/grades", AuthMiddleware, ListGrades)
	app.Get("/grades/statistics", AuthMiddleware, GetGradeStatistics)
	app.Get("/grades/:id", AuthMiddleware, GetGrade)
	app.Post("/grades", AuthMiddleware, CreateGrade)
	app.Patch("/grades/:id", AuthMiddleware, UpdateGrade)
	app.Post("/grades/batch", AuthMiddleware, BatchCreateGrades)

	app.Get("/certificates", AuthMiddleware, GenerateCertificate)

	log.Info("REST API started")
}
