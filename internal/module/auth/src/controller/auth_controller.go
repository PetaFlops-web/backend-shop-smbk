package controller

import (
	"github.com/PetaFlops-web/backend-shop-smbk/internal/module/auth/src/model"
	"github.com/PetaFlops-web/backend-shop-smbk/internal/module/auth/src/usecase"
	"github.com/PetaFlops-web/backend-shop-smbk/internal/shared/middleware"
	"github.com/PetaFlops-web/backend-shop-smbk/internal/shared/response"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type AuthController struct {
	Log     *logrus.Logger
	UseCase *usecase.AuthUseCase
}

func NewAuthController(useCase *usecase.AuthUseCase, logger *logrus.Logger) *AuthController {
	return &AuthController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *AuthController) Current(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	resp, err := c.UseCase.Current(ctx.UserContext(), auth.ID)
	if err != nil {
		c.Log.Warnf("Failed to get current user : %+v", err)
		return err
	}

	return ctx.JSON(response.WebResponse[*model.UserResponse]{
		Data:    resp,
		Message: "Get current user successful",
		Success: true,
	})
}

func (c *AuthController) Register(ctx *fiber.Ctx) error {
	request := new(model.RegisterRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, "Format data request tidak valid")
	}

	resp, err := c.UseCase.Register(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to register user : %+v", err)
		return err
	}

	return ctx.JSON(response.WebResponse[*model.AuthResponse]{
		Data:    resp,
		Message: "Register successful",
		Success: true,
	})
}

func (c *AuthController) Login(ctx *fiber.Ctx) error {
	request := new(model.LoginRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, "Format data request tidak valid")
	}

	resp, err := c.UseCase.Login(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to login user : %+v", err)
		return err
	}

	return ctx.JSON(response.WebResponse[*model.AuthResponse]{
		Data:    resp,
		Message: "Login successful",
		Success: true,
	})
}