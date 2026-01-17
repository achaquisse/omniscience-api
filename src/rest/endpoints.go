package rest

import (
	"github.com/gofiber/fiber/v2"
)

func Init(app *fiber.App) {
	SetupSwagger(app)

	app.Get("/student-classes", AuthMiddleware, ListStudentClass)
	app.Get("/registrations", AuthMiddleware, ListRegistrations)
	app.Post("/attendance", AuthMiddleware, RecordAttendance)
	app.Post("/attendance/bulk", AuthMiddleware, RecordBulkAttendance)
	app.Get("/attendance/report", AuthMiddleware, GetStudentAttendanceReport)
	app.Get("/attendance/class-report", AuthMiddleware, GetClassAttendanceReport)
}
