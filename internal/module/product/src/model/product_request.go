package model

type SearchProductRequest struct {
	StoreId string `json:"store_id" query:"store_id"`
	Name    string `json:"name" query:"name"`
	Page    int    `json:"page" query:"page"`
	Size    int    `json:"size" query:"size"`
}

type CreateProductRequest struct {
	StoreID      string `json:"store_id" validate:"required"`
	ProductName  string `json:"product_name" validate:"required"`
	CostPrice    int64  `json:"cost_price" validate:"min=0"`
	SellingPrice int64  `json:"selling_price" validate:"required,min=0"`
	Stock        int    `json:"stock" validate:"min=0"`
	Unit         string `json:"unit" validate:"required"`
}

type UpdateProductRequest struct {
	ID           string `json:"-" validate:"required"` // Dari path param
	StoreID      string `json:"store_id" validate:"required"`
	ProductName  string `json:"product_name" validate:"required"`
	CostPrice    int64  `json:"cost_price" validate:"min=0"`
	SellingPrice int64  `json:"selling_price" validate:"required,min=0"`
	Stock        int    `json:"stock" validate:"min=0"`
	Unit         string `json:"unit" validate:"required"`
}