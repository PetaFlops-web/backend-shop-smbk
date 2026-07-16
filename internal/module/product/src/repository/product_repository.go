package repository

import (
	"github.com/PetaFlops-web/backend-shop-smbk/internal/module/product/src/entity"
	"github.com/PetaFlops-web/backend-shop-smbk/internal/module/product/src/model"
	"github.com/PetaFlops-web/backend-shop-smbk/internal/shared/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ProductRepository struct {
	repository.Repository[entity.Product]
	Log *logrus.Logger
}

func NewProductRepository(log *logrus.Logger) *ProductRepository {
	return &ProductRepository{Log: log}
}

func (r *ProductRepository) FindByStoreId(db *gorm.DB, storeId string) ([]entity.Product, error) {
	var products []entity.Product
	if err := db.Where("store_id = ?", storeId).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepository) Search(db *gorm.DB, request *model.SearchProductRequest) ([]entity.Product, int64, error) {
	var products []entity.Product
	if err := db.Scopes(r.FilterProduct(request)).Offset((request.Page - 1) * request.Size).Limit(request.Size).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	var total int64 = 0
	if err := db.Model(&entity.Product{}).Scopes(r.FilterProduct(request)).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *ProductRepository) FilterProduct(request *model.SearchProductRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		tx = tx.Where("store_id = ?", request.StoreId)

		if name := request.Name; name != "" {
			name = "%" + name + "%"
			tx = tx.Where("product_name LIKE ?", name)
		}

		return tx
	}
}