package server

import (
	"monitoring/config"
	"monitoring/internal/domain/handler"
	"monitoring/internal/domain/repository"
	"monitoring/internal/domain/usecase"
	"monitoring/internal/middleware"
	"monitoring/pkg/db"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"gorm.io/gorm"

	database "monitoring/pkg/db"
)

func SetupRoutes(app *fiber.App, cfg *config.Config, db *gorm.DB, influx influxdb2.Client, redis0 *db.Client) {

	monitoringRepo := repository.NewMonitoringRepository(influx, &cfg.InfluxDB)

	mqttClient := database.NewMQTTClient(cfg, monitoringRepo)
	go mqttClient.Start()
	var validate = validator.New()

	autRepo := repository.NewUserRepository(db)
	authUsecase := usecase.NewAuthUsecase(autRepo, redis0)
	authHandler := handler.NewAuthHandler(authUsecase, validate)

	devRepo := repository.NewDeviceRepository(db)
	devUsecase := usecase.NewDeviceUsecase(devRepo, redis0)
	deviceHandler := handler.NewDeviceHandler(devUsecase, validate)

	telemetryHandler := handler.NewMQTTHandler(mqttClient, validate)
	monitoringUsecase := usecase.NewMonitoringUsecase(monitoringRepo)
	monitoringHandler := handler.NewMonitoringHandler(monitoringUsecase, validate)

	api := app.Group("/api/v1")
	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.RefreshToken)

	// Protected routes
	protected := api.Use(middleware.JWTMiddleware())

	// // Device routes
	devices := protected.Group("/devices")
	devices.Post("/", deviceHandler.CreateDevice)
	devices.Get("/", deviceHandler.GetDevice)
	// devices.Get("/:id", deviceHandler.GetDevice)
	// devices.Put("/:id", deviceHandler.UpdateDevice)
	// devices.Delete("/:id", deviceHandler.DeleteDevice)

	// // Telemetry routes
	telemetry := protected.Group("/telemetry")
	telemetry.Get("/device/:device_id", monitoringHandler.GetMonitoringByDeviceID)
	// telemetry.Get("/device/:deviceId/latest", telemetryHandler.GetLatestTelemetry)
	telemetry.Post("/:id", telemetryHandler.TriggerMQTT)
}
