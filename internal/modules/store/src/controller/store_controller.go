package controller

import (
	"github.com/PetaFlops-web/backend-shop-smbk/internal/modules/store/src/model"
	"github.com/PetaFlops-web/backend-shop-smbk/internal/modules/store/src/usecase"
	"github.com/PetaFlops-web/backend-shop-smbk/internal/shared/middleware"
	"github.com/PetaFlops-web/backend-shop-smbk/internal/shared/response"
	"github.com/gofiber/fiber/v2"
)

type StoreController struct {
	UseCase usecase.StoreUseCase
}

func NewStoreController(useCase usecase.StoreUseCase) *StoreController {
	return &StoreController{
		UseCase: useCase,
	}
}

// Create godoc
// @Summary      Create a new store
// @Description  Creates a new store for the currently authenticated user
// @Tags         Store
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body model.CreateStoreRequest true "Store details"
// @Success      200  {object}  response.WebResponse[model.StoreResponse]
// @Failure      400  {object}  response.ApiErrorResponse
// @Failure      401  {object}  response.ApiErrorResponse
// @Failure      409  {object}  response.ApiErrorResponse
// @Router       /stores [post]
func (c *StoreController) Create(ctx *fiber.Ctx) error {
	var request model.CreateStoreRequest
	if err := ctx.BodyParser(&request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	auth := middleware.GetUser(ctx)
	resp, err := c.UseCase.Create(ctx.UserContext(), auth.ID, &request)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.WebResponse[*model.StoreResponse]{
		Data: resp,
	})
}

// GetMyStore godoc
// @Summary      Get my store
// @Description  Returns the store of the currently authenticated user
// @Tags         Store
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.WebResponse[model.StoreResponse]
// @Failure      401  {object}  response.ApiErrorResponse
// @Failure      404  {object}  response.ApiErrorResponse
// @Router       /stores [get]
func (c *StoreController) GetMyStore(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	resp, err := c.UseCase.GetMyStore(ctx.UserContext(), auth.ID)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.WebResponse[*model.StoreResponse]{
		Data: resp,
	})
}

// UpdateMyStore godoc
// @Summary      Update my store
// @Description  Updates the store of the currently authenticated user
// @Tags         Store
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body model.UpdateStoreRequest true "Store update details"
// @Success      200  {object}  response.WebResponse[model.StoreResponse]
// @Failure      400  {object}  response.ApiErrorResponse
// @Failure      401  {object}  response.ApiErrorResponse
// @Failure      404  {object}  response.ApiErrorResponse
// @Router       /stores [put]
func (c *StoreController) UpdateMyStore(ctx *fiber.Ctx) error {
	var request model.UpdateStoreRequest
	if err := ctx.BodyParser(&request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	auth := middleware.GetUser(ctx)
	resp, err := c.UseCase.UpdateMyStore(ctx.UserContext(), auth.ID, &request)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.WebResponse[*model.StoreResponse]{
		Data: resp,
	})
}
