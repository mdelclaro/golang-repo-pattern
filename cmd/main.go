package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"

	"golang-repo-pattern/internal/infra/database"
)

func main() {
	database.StartDb()
	app := fiber.New()

	app.Use(cors.New())
	app.Use(func(c fiber.Ctx) error {
		return fiber.ErrNotFound
	})

	log.Fatal(app.Listen(os.Getenv("APP_PORT")))
}
