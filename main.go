package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})

	// separate logfile for app
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err == nil {
		log.SetOutput(io.MultiWriter(os.Stdout, file))
	} else {
		log.WithError(err).Error("Failed to log to file, using default stdout")
		log.SetOutput(os.Stdout)
	}
}

func main() {

	// separate logfile for fiber
	file, err := os.OpenFile("fiber.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	webApp := fiber.New()
	webApp.Use(logger.New(logger.Config{
		Format:     "[${time}] [${ip}]:${port} ${status} - ${latency} ${method} ${path}\n",
		TimeFormat: "02-Jan-2006 15:04:05",
		TimeZone:   "Local",
		Output:     file,
	}))

	webApp.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	webApp.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	webApp.Get("/pages", func(c *fiber.Ctx) error {
		pageNumber := c.Query("page")
		toSort := c.Query("sort")
		return c.SendString("Page called: " + pageNumber + "\n" + "Sort: " + toSort)
	})

	port := "8088"

	log.WithFields(log.Fields{
		"port": port,
	}).Info("Server started")

	log.Fatal(webApp.Listen(":" + port))
}
