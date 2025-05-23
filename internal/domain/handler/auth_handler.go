package handler

import (
	"monitoring/internal/domain/entity"
	iface "monitoring/internal/domain/interface"
	"monitoring/pkg/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authUsecase iface.AuthUsecase
	validate    *validator.Validate
}

func NewAuthHandler(au iface.AuthUsecase, validate *validator.Validate) *AuthHandler {
	return &AuthHandler{
		authUsecase: au,
		validate:    validate,
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	req := new(entity.RegisterRequest)
	if err := c.BodyParser(req); err != nil {
		return fiber.ErrBadRequest
	}

	if err := h.validate.Struct(req); err != nil {
		validationErrors := utils.FormatValidationErrors(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": validationErrors,
		})
	}

	err := h.authUsecase.Register(c.Context(), req)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "user registered successfully",
	})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	req := new(entity.LoginRequest)
	if err := c.BodyParser(req); err != nil {
		return fiber.ErrBadRequest
	}

	accessToken, refreshToken, err := h.authUsecase.Login(c.Context(), req.Username, req.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	return c.JSON(fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	req := new(entity.RefreshTokenRequest)
	if err := c.BodyParser(req); err != nil {
		return fiber.ErrBadRequest
	}

	newAccessToken, newRefreshToken, err := h.authUsecase.RefreshToken(c.Context(), req.RefreshToken)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	return c.JSON(fiber.Map{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})
}
