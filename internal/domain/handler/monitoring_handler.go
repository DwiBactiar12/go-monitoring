package handler

import (
	"fmt"
	"monitoring/internal/domain/entity"
	iface "monitoring/internal/domain/interface"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type MonitoringHandler struct {
	usecase  iface.MonitoringUseCase
	validate *validator.Validate
}

func NewMonitoringHandler(uc iface.MonitoringUseCase, validate *validator.Validate) *MonitoringHandler {
	return &MonitoringHandler{
		usecase:  uc,
		validate: validate,
	}
}

// POST /monitoring
// Simpan data monitoring
func (h *MonitoringHandler) StoreMonitoring(ctx *fiber.Ctx) error {
	var data entity.MonitoringData

	if err := ctx.BodyParser(&data); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	// Validasi struct
	if err := h.validate.Struct(&data); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Simpan data
	if err := h.usecase.StoreMonitoringData(ctx.Context(), &data); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(fiber.Map{"message": "Data monitoring berhasil disimpan"})
}

// GET /monitoring/:device_id?start=2023-01-01T00:00:00Z&end=2023-01-02T00:00:00Z&limit=100
// Ambil data monitoring berdasarkan device_id dan range waktu
func (h *MonitoringHandler) GetMonitoringByDeviceID(ctx *fiber.Ctx) error {
	deviceIDStr := ctx.Params("device_id")
	fmt.Println(deviceIDStr)
	deviceID, err := uuid.Parse(deviceIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid device_id"})
	}

	startStr := ctx.Query("start")
	endStr := ctx.Query("end")
	limitStr := ctx.Query("limit", "100")

	startTime, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid start time"})
	}
	endTime, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid end time"})
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid limit"})
	}

	data, err := h.usecase.GetMonitoringDataByDevice(ctx.Context(), deviceID, startTime, endTime, limit)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(data)
}
