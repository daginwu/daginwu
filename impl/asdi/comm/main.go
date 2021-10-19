package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type App struct {
	server *fiber.App
	db     *Database
}

func main() {

	// Init application
	app := &App{
		server: fiber.New(),
		db:     InitDatabase(),
	}
	// Logger
	app.server.Use(logger.New())
	// Init HTTP server

	app.server.Get("/hello", Hello)

	// User management
	app.server.Get("/api/users", app.GetUsers)
	app.server.Get("/api/users/:id", app.GetUser)
	app.server.Post("/api/users", app.CreateUser)
	app.server.Delete("/api/users/:id", app.DeleteUser)

	// Txn management
	app.server.Post("/api/txns", app.CreateTxn)
	app.server.Listen(":3000")
}

func Hello(c *fiber.Ctx) error {
	return c.SendString("Hello Google DSC in NTHU")
}
