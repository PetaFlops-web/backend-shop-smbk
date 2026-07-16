package model

type ProductResponse struct {
	ID           string `json:"id"`
	StoreID      string `json:"store_id"`
	ProductName  string `json:"product_name"`
	CostPrice    int64  `json:"cost_price"`
	SellingPrice int64  `json:"selling_price"`
	Stock        int    `json:"stock"`
	Unit         string `json:"unit"`
	CreatedAt    int64  `json:"created_at"`
	UpdatedAt    int64  `json:"updated_at"`
}