package main

import (
	"skulla-api/db"
	"skulla-api/rest"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Connects to database server
	db.Connect()

	// Instantiate web server
	app := fiber.New()
	rest.Init(app) //Init rest endpoints
	err := app.Listen(":8080")
	if err != nil {
		return
	}
}
