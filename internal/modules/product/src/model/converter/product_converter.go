package converter

import (
	"github.com/PetaFlops-web/backend-shop-smbk/internal/modules/product/src/entity"
	"github.com/PetaFlops-web/backend-shop-smbk/internal/modules/product/src/model"
)

func ProductToResponse(product *entity.Product) *model.ProductResponse {
	if product == nil {	
		return nil
	}
	return &model.ProductResponse{
		ID:           product.ID,
		StoreID:      product.StoreID,
		ProductName:  product.ProductName,
		CostPrice:    product.CostPrice,
		SellingPrice: product.SellingPrice,
		Stock:        product.Stock,
		Unit:         product.Unit,
		CreatedAt:    product.CreatedAt,
		UpdatedAt:    product.UpdatedAt,
	}
}