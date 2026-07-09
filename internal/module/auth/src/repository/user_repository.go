package repository

import (
	"github.com/PetaFlops-web/backend-shop-smbk/internal/module/auth/src/entity"
	"github.com/PetaFlops-web/backend-shop-smbk/internal/shared/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepository struct {
	repository.Repository[entity.User]
	Log *logrus.Logger
}

func NewUserRepository(log *logrus.Logger) *UserRepository {
	return &UserRepository{Log: log}
}

func (r *UserRepository) FindByUsername(db *gorm.DB, user *entity.User, username string) error {
	return db.Where("username = ?", username).Take(user).Error
}

func (r *UserRepository) CountByUsernameOrEmail(db *gorm.DB, username string, email string) (int64, error) {
	var count int64
	err := db.Model(&entity.User{}).Where("username = ? OR email = ?", username, email).Count(&count).Error
	return count, err
}