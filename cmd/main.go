package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"

	"golang-repo-pattern/internal/config"
	"golang-repo-pattern/internal/domain/device"
	"golang-repo-pattern/internal/infra/database"
)

func main() {
	database.StartDb()
	app := fiber.New()
	app.Use(cors.New())

	deviceRepository := device.NewRepository(database.DB.Db)
	deviceService := device.NewService(device.ServiceParams{
		Repo: deviceRepository,
	})

	device.NewHttpHandler(app, deviceService)

	app.Use(func(c fiber.Ctx) error {
		return fiber.ErrNotFound
	})

	log.Fatal(app.Listen(config.GetEnv("APP_PORT")))
}
