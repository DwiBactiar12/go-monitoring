package handler

import (
	"encoding/json"
	"fmt"
	"monitoring/internal/domain/entity"
	"monitoring/pkg/db"
	"monitoring/pkg/utils"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type MQTTHandler struct {
	mqtt     *db.MQTTClient
	validate *validator.Validate
}

func NewMQTTHandler(m *db.MQTTClient, validate *validator.Validate) *MQTTHandler {
	return &MQTTHandler{
		mqtt:     m,
		validate: validate,
	}
}

func (m *MQTTHandler) TriggerMQTT(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing device ID in URL",
		})
	}

	var payload entity.MonitoringRequest
	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	if err := m.validate.Struct(payload); err != nil {
		validationErrors := utils.FormatValidationErrors(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": validationErrors,
		})
	}

	sendData := entity.MonitoringData{
		DeviceID:    uuid.MustParse(id),
		CPUUsage:    payload.CPUUsage,
		MemoryUsage: payload.MemoryUsage,
		DiskUsage:   payload.DiskUsage,
		Temperature: payload.Temperature,
		Timestamp:   time.Now(),
	}

	topic := fmt.Sprintf("iot/monitoring/%s/telemetry", id)
	payloadBytes, _ := json.Marshal(sendData)

	if err := m.mqtt.PublishTelemetry(topic, payloadBytes); err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(fiber.Map{"message": "Published to MQTT"})
}
