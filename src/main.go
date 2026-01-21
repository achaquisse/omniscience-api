package main

import (
	"skulla-api/db"
	"skulla-api/rest"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// Connects to database server
	db.Connect()

	// Instantiate web server
	app := fiber.New()

	// Configure CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	rest.Init(app) //Init rest endpoints
	err := app.Listen(":8080")
	if err != nil {
		return
	}
}
