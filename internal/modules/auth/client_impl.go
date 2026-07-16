package auth

import (
	"context"

	auth_client "github.com/PetaFlops-web/backend-shop-smbk/internal/modules/auth-client"
	"github.com/PetaFlops-web/backend-shop-smbk/internal/modules/auth/src/entity"
	"gorm.io/gorm"
)

type clientImpl struct {
	db *gorm.DB
}

func (c *clientImpl) GetUserByID(ctx context.Context, userID string) (*auth_client.UserDTO, error) {
	var user entity.User
	err := c.db.WithContext(ctx).Where("id = ?", userID).Take(&user).Error
	if err != nil {
		return nil, err
	}

	return &auth_client.UserDTO{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}