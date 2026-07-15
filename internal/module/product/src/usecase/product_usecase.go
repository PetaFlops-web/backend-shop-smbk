package usecase

import (
	"context"

	"github.com/PetaFlops-web/backend-shop-smbk/internal/module/product/src/entity"
	"github.com/PetaFlops-web/backend-shop-smbk/internal/module/product/src/model"
	"github.com/PetaFlops-web/backend-shop-smbk/internal/module/product/src/model/converter"
	"github.com/PetaFlops-web/backend-shop-smbk/internal/module/product/src/repository"
	"github.com/PetaFlops-web/backend-shop-smbk/internal/shared/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ProductUseCase struct {
	DB                *gorm.DB
	Log               *logrus.Logger
	Validate          *validator.Validate
	ProductRepository *repository.ProductRepository
}

func NewProductUseCase(
	db *gorm.DB,
	log *logrus.Logger,
	validate *validator.Validate,
	productRepo *repository.ProductRepository,
) *ProductUseCase {
	return &ProductUseCase{
		DB:                db,
		Log:               log,
		Validate:          validate,
		ProductRepository: productRepo,
	}
}

func (c *ProductUseCase) Create(ctx context.Context, request *model.CreateProductRequest) (*model.ProductResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "Format data request tidak valid")
	}

	// Generate Product ID
	productID, err := utils.GenerateProductId()
	if err != nil {
		c.Log.Warnf("Failed to generate product id : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Gagal membuat ID produk")
	}

	product := &entity.Product{
		ID:           productID,
		StoreID:      request.StoreID,
		ProductName:  request.ProductName,
		CostPrice:    request.CostPrice,
		SellingPrice: request.SellingPrice,
		Stock:        request.Stock,
		Unit:         request.Unit,
	}

	if err := c.ProductRepository.Create(tx, product); err != nil {
		c.Log.Warnf("Failed create product : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Gagal menyimpan data produk")
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Gagal menyimpan data produk")
	}

	return converter.ProductToResponse(product), nil
}

func (c *ProductUseCase) Update(ctx context.Context, request *model.UpdateProductRequest) (*model.ProductResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "Format data request tidak valid")
	}

	product := new(entity.Product)
	if err := c.ProductRepository.FindById(tx, product, request.ID); err != nil {
		c.Log.Warnf("Product not found : %+v", err)
		return nil, fiber.NewError(fiber.StatusNotFound, "Produk tidak ditemukan")
	}

	if product.StoreID != request.StoreID {
		c.Log.Warnf("Forbidden: Store ID mismatch. Expected %s, got %s", product.StoreID, request.StoreID)
		return nil, fiber.NewError(fiber.StatusForbidden, "Akses ditolak")
	}

	product.ProductName = request.ProductName
	product.CostPrice = request.CostPrice
	product.SellingPrice = request.SellingPrice
	product.Stock = request.Stock
	product.Unit = request.Unit

	if err := c.ProductRepository.Update(tx, product); err != nil {
		c.Log.Warnf("Failed update product : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Gagal mengubah data produk")
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Gagal mengubah data produk")
	}

	return converter.ProductToResponse(product), nil
}

func (c *ProductUseCase) Get(ctx context.Context, storeId string, id string) (*model.ProductResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	product := new(entity.Product)
	if err := c.ProductRepository.FindById(tx, product, id); err != nil {
		c.Log.Warnf("Product not found : %+v", err)
		return nil, fiber.NewError(fiber.StatusNotFound, "Produk tidak ditemukan")
	}

	if product.StoreID != storeId {
		c.Log.Warnf("Forbidden: Store ID mismatch. Expected %s, got %s", product.StoreID, storeId)
		return nil, fiber.NewError(fiber.StatusForbidden, "Akses ditolak")
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Gagal mengambil data produk")
	}

	return converter.ProductToResponse(product), nil
}

func (c *ProductUseCase) Delete(ctx context.Context, storeId string, id string) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	product := new(entity.Product)
	if err := c.ProductRepository.FindById(tx, product, id); err != nil {
		c.Log.Warnf("Product not found : %+v", err)
		return fiber.NewError(fiber.StatusNotFound, "Produk tidak ditemukan")
	}

	if product.StoreID != storeId {
		c.Log.Warnf("Forbidden: Store ID mismatch. Expected %s, got %s", product.StoreID, storeId)
		return fiber.NewError(fiber.StatusForbidden, "Akses ditolak")
	}

	if err := c.ProductRepository.Delete(tx, product); err != nil {
		c.Log.Warnf("Failed delete product : %+v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal menghapus produk")
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal menghapus produk")
	}

	return nil
}

func (c *ProductUseCase) Search(ctx context.Context, request *model.SearchProductRequest) ([]model.ProductResponse, int64, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request query : %+v", err)
		return nil, 0, fiber.NewError(fiber.StatusBadRequest, "Format pencarian tidak valid")
	}

	if request.StoreId == "" {
		return nil, 0, fiber.NewError(fiber.StatusBadRequest, "Store ID wajib diisi")
	}

	products, total, err := c.ProductRepository.Search(tx, request)
	if err != nil {
		c.Log.Warnf("Failed search products : %+v", err)
		return nil, 0, fiber.NewError(fiber.StatusInternalServerError, "Gagal mencari produk")
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, 0, fiber.NewError(fiber.StatusInternalServerError, "Gagal mencari produk")
	}

	responses := make([]model.ProductResponse, len(products))
	for i, product := range products {
		responses[i] = *converter.ProductToResponse(&product)
	}

	return responses, total, nil
}