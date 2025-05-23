package main

import (
	"log"
	"monitoring/config"
	"monitoring/internal/server"
	database "monitoring/pkg/db"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize InfluxDB
	influxClient, err := database.NewInfluxDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to InfluxDB:", err)
	}

	redis0, err := database.NewRedisClient(&cfg.Redis, 0)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize MQTT client
	

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())

	// Routes
	server.SetupRoutes(app, cfg, db, influxClient, redis0)

	// Start server
	log.Printf("Server starting on %s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Fatal(app.Listen(cfg.Server.Host + ":" + cfg.Server.Port))
}
