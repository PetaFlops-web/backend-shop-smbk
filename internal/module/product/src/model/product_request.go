package model

type SearchProductRequest struct {
	StoreId string `json:"store_id"`
	Name    string `json:"name"`
	Page    int    `json:"page"`
	Size    int    `json:"size"`
}