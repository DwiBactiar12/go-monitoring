package handler

import (
	"monitoring/internal/domain/entity"
	iface "monitoring/internal/domain/interface"
	"monitoring/pkg/jwt"
	"monitoring/pkg/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type deviceHandler struct {
	deviceUsecase iface.DeviceUseCase
	validate      *validator.Validate
}

func NewDeviceHandler(deviceUsecase iface.DeviceUseCase, validate *validator.Validate) *deviceHandler {
	return &deviceHandler{
		deviceUsecase: deviceUsecase,
		validate:      validate,
	}
}

func (d *deviceHandler) CreateDevice(ctx *fiber.Ctx) error {
	var reqDevice entity.DeviceRequest
	if err := ctx.BodyParser(&reqDevice); err != nil {
		return err
	}

	if err := d.validate.Struct(reqDevice); err != nil {
		validationErrors := utils.FormatValidationErrors(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": validationErrors,
		})
	}
	claims, ok := ctx.Locals("user").(*jwt.Claims)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user data"})
	}

	err := d.deviceUsecase.Create(ctx.Context(), &entity.Device{
		Name:       reqDevice.Name,
		Type:       reqDevice.Type,
		MacAddress: reqDevice.MacAddress,
		IPAddress:  reqDevice.IPAddress,
		Location:   reqDevice.Location,
		UserID:     uuid.MustParse(claims.UserID),
	})

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Device registered successfully",
	})
}

func (d *deviceHandler) GetDevice(ctx *fiber.Ctx) error {
	_, ok := ctx.Locals("user").(*jwt.Claims)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user data"})
	}

	devices, err := d.deviceUsecase.List(ctx.Context(), 10, 0)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": devices,
	})
}
