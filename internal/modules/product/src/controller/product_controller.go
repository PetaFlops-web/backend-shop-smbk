package controller

import (
	"math"

	"github.com/PetaFlops-web/backend-shop-smbk/internal/module/product/src/model"
	"github.com/PetaFlops-web/backend-shop-smbk/internal/module/product/src/usecase"
	"github.com/PetaFlops-web/backend-shop-smbk/internal/shared/response"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type ProductController struct {
	Log     *logrus.Logger
	UseCase *usecase.ProductUseCase
}

func NewProductController(useCase *usecase.ProductUseCase, logger *logrus.Logger) *ProductController {
	return &ProductController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *ProductController) Create(ctx *fiber.Ctx) error {
	request := new(model.CreateProductRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, "Format data request tidak valid")
	}

	resp, err := c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to create product : %+v", err)
		return err
	}

	return ctx.JSON(response.WebResponse[*model.ProductResponse]{
		Data:    resp,
		Message: "Berhasil menambahkan produk",
		Success: true,
	})
}

func (c *ProductController) Update(ctx *fiber.Ctx) error {
	request := new(model.UpdateProductRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, "Format data request tidak valid")
	}

	request.ID = ctx.Params("id")

	resp, err := c.UseCase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to update product : %+v", err)
		return err
	}

	return ctx.JSON(response.WebResponse[*model.ProductResponse]{
		Data:    resp,
		Message: "Berhasil mengubah data produk",
		Success: true,
	})
}

func (c *ProductController) Get(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	storeId := ctx.Query("store_id")

	if storeId == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Store ID tidak boleh kosong")
	}

	resp, err := c.UseCase.Get(ctx.UserContext(), storeId, id)
	if err != nil {
		c.Log.Warnf("Failed to get product : %+v", err)
		return err
	}

	return ctx.JSON(response.WebResponse[*model.ProductResponse]{
		Data:    resp,
		Message: "Berhasil mengambil data produk",
		Success: true,
	})
}

func (c *ProductController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	storeId := ctx.Query("store_id")

	if storeId == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Store ID tidak boleh kosong")
	}

	err := c.UseCase.Delete(ctx.UserContext(), storeId, id)
	if err != nil {
		c.Log.Warnf("Failed to delete product : %+v", err)
		return err
	}

	return ctx.JSON(response.WebResponse[any]{
		Data:    nil,
		Message: "Berhasil menghapus produk",
		Success: true,
	})
}

func (c *ProductController) List(ctx *fiber.Ctx) error {
	request := &model.SearchProductRequest{
		StoreId: ctx.Query("store_id", ""),
		Name:    ctx.Query("name", ""),
		Page:    ctx.QueryInt("page", 1),
		Size:    ctx.QueryInt("size", 10),
	}

	responses, total, err := c.UseCase.Search(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed searching product : %+v", err)
		return err
	}

	paging := &response.PageMetadata{
		Page:      request.Page,
		Size:      request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
	}

	return ctx.JSON(response.WebResponse[[]model.ProductResponse]{
		Data:    responses,
		Paging:  paging,
		Message: "Berhasil menampilkan daftar produk",
		Success: true,
	})
}