package model

type StoreResponse struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	StoreName string `json:"store_name"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}
