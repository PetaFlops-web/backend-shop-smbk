package controller

import (
	"github.com/PetaFlops-web/backend-shop-smbk/internal/modules/auth/src/model"
	"github.com/PetaFlops-web/backend-shop-smbk/internal/modules/auth/src/usecase"
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

// Current godoc
// @Summary      Get current logged-in user
// @Description  Returns the profile of the currently authenticated user
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.WebResponse[model.UserResponse]
// @Failure      401  {object}  response.ApiErrorResponse
// @Router       /users/_current [get]
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

// Register godoc
// @Summary      Register a new merchant user
// @Description  Creates a new user with username and password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      model.RegisterRequest  true  "Register Request"
// @Success      200   {object}  response.WebResponse[model.AuthResponse]
// @Failure      400   {object}  response.ApiErrorResponse
// @Router       /users [post]
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

// Login godoc
// @Summary      Login merchant user
// @Description  Authenticates a user and returns a JWT token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      model.LoginRequest  true  "Login Request"
// @Success      200   {object}  response.WebResponse[model.AuthResponse]
// @Failure      400   {object}  response.ApiErrorResponse
// @Failure      401   {object}  response.ApiErrorResponse
// @Router       /users/_login [post]
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