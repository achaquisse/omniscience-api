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

	app.Get("/evaluation-categories", AuthMiddleware, ListEvaluationCategories)
	app.Get("/evaluation-categories/:id", AuthMiddleware, GetEvaluationCategory)
	app.Post("/evaluation-categories", AuthMiddleware, CreateEvaluationCategory)
	app.Put("/evaluation-categories/:id", AuthMiddleware, UpdateEvaluationCategory)
	app.Delete("/evaluation-categories/:id", AuthMiddleware, DeleteEvaluationCategory)

	app.Get("/grading-formulas", AuthMiddleware, GetGradingFormula)
	app.Post("/grading-formulas", AuthMiddleware, CreateGradingFormula)
	app.Put("/grading-formulas/:id", AuthMiddleware, UpdateGradingFormula)

	app.Get("/evaluation-items", AuthMiddleware, ListEvaluationItems)
	app.Get("/evaluation-items/:id", AuthMiddleware, GetEvaluationItem)
	app.Post("/evaluation-items", AuthMiddleware, CreateEvaluationItem)
	app.Put("/evaluation-items/:id", AuthMiddleware, UpdateEvaluationItem)
	app.Delete("/evaluation-items/:id", AuthMiddleware, DeleteEvaluationItem)

	app.Get("/grades", AuthMiddleware, ListGrades)
	app.Get("/grades/:id", AuthMiddleware, GetGrade)
	app.Post("/grades", AuthMiddleware, CreateGrade)
	app.Put("/grades/:id", AuthMiddleware, UpdateGrade)
	app.Post("/grades/batch", AuthMiddleware, BatchCreateGrades)
	app.Post("/grades/publish/:evaluation_item_id", AuthMiddleware, PublishGrades)

	app.Get("/reports/student/:registration_id", AuthMiddleware, GetStudentReport)
	app.Get("/reports/class/:student_class_id", AuthMiddleware, GetClassReport)
	app.Get("/final-grades/:registration_id", AuthMiddleware, GetFinalGrade)
	app.Post("/final-grades/calculate/:student_class_id", AuthMiddleware, CalculateFinalGrades)

	log.Info("REST API started")
}
