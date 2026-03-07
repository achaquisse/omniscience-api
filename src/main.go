package main

import (
	"skulla-api/db"
	"skulla-api/rest"
	"skulla-api/scheduler"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	db.Connect()

	scheduler.InitAttendanceScheduler()
	defer scheduler.StopScheduler()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, PATCH, POST, PUT, DELETE, OPTIONS",
	}))

	rest.Init(app)
	err := app.Listen(":8080")
	if err != nil {
		return
	}
}
