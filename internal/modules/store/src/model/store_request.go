package model

type StoreRequest struct {
	ID           string `json:"-"`
	UserID       string `json:"-" validate:"required"`
	StoreName    string `json:"store_name" validate:"required,min=3"`
}

