package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func main() {
	app := fiber.New()

	app.Use(cors.New())
	app.Use(func(c fiber.Ctx) error {
		return fiber.ErrNotFound
	})

	log.Fatal(app.Listen(os.Getenv("APP_PORT")))
}
