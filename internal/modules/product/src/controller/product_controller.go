package controller

import (
	"math"

	"github.com/PetaFlops-web/backend-shop-smbk/internal/modules/product/src/model"
	"github.com/PetaFlops-web/backend-shop-smbk/internal/modules/product/src/usecase"
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

// Create godoc
// @Summary      Menambahkan data produk baru
// @Description  Membuat entri produk baru untuk toko tertentu
// @Tags         Product
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      model.CreateProductRequest  true  "Data Produk"
// @Success      200   {object}  response.WebResponse[model.ProductResponse]
// @Failure      400   {object}  response.ApiErrorResponse
// @Failure      500   {object}  response.ApiErrorResponse
// @Router       /products [post]
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

// Update godoc
// @Summary      Mengubah data produk yang sudah ada
// @Description  Mengubah detail produk berdasarkan ID
// @Tags         Product
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string                     true  "Product ID"
// @Param        body  body      model.UpdateProductRequest true  "Data Produk"
// @Success      200   {object}  response.WebResponse[model.ProductResponse]
// @Failure      400   {object}  response.ApiErrorResponse
// @Failure      500   {object}  response.ApiErrorResponse
// @Router       /products/{id} [put]
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

// Get godoc
// @Summary      Menampilkan detail satu produk
// @Description  Mengambil data produk berdasarkan ID dan Store ID
// @Tags         Product
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string  true  "Product ID"
// @Param        store_id query     string  true  "Store ID"
// @Success      200      {object}  response.WebResponse[model.ProductResponse]
// @Failure      400      {object}  response.ApiErrorResponse
// @Failure      404      {object}  response.ApiErrorResponse
// @Router       /products/{id} [get]
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

// Delete godoc
// @Summary      Menghapus produk
// @Description  Menghapus data produk berdasarkan ID dan Store ID
// @Tags         Product
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string  true  "Product ID"
// @Param        store_id query     string  true  "Store ID"
// @Success      200      {object}  response.WebResponse[any]
// @Failure      400      {object}  response.ApiErrorResponse
// @Failure      404      {object}  response.ApiErrorResponse
// @Router       /products/{id} [delete]
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

// List godoc
// @Summary      Menampilkan daftar produk dengan pagination
// @Description  Mencari dan menampilkan daftar produk untuk suatu toko
// @Tags         Product
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        store_id query     string  true   "Store ID"
// @Param        name     query     string  false  "Keyword nama produk"
// @Param        page     query     int     false  "Nomor Halaman" default(1)
// @Param        size     query     int     false  "Ukuran Halaman" default(10)
// @Success      200      {object}  response.WebResponse[[]model.ProductResponse]
// @Failure      400      {object}  response.ApiErrorResponse
// @Failure      500      {object}  response.ApiErrorResponse
// @Router       /products [get]
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