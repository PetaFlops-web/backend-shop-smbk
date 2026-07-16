package auth_client

import "context"

type UserDTO struct {
	ID       string
	Username string
	Email    string
}

type Client interface {
	GetUserByID(ctx context.Context, userID string) (*UserDTO, error)
}