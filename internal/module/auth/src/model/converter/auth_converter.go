package converter

import (
	"github.com/PetaFlops-web/backend-shop-smbk/internal/module/auth/src/entity"
	"github.com/PetaFlops-web/backend-shop-smbk/internal/module/auth/src/model"
)

func UserToResponse(user *entity.User) *model.UserResponse {
	return &model.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}