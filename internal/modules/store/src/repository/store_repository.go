package repository

import (
	"context"

	"github.com/PetaFlops-web/backend-shop-smbk/internal/modules/store/src/entity"
	"gorm.io/gorm"
)

type StoreRepository interface {
	Create(ctx context.Context, store *entity.Store) error
	FindByUserID(ctx context.Context, userID string) (*entity.Store, error)
	FindByID(ctx context.Context, id string) (*entity.Store, error)
	Update(ctx context.Context, store *entity.Store) error
}

type storeRepositoryImpl struct {
	db *gorm.DB
}

func NewStoreRepository(db *gorm.DB) StoreRepository {
	return &storeRepositoryImpl{
		db: db,
	}
}

func (r *storeRepositoryImpl) Create(ctx context.Context, store *entity.Store) error {
	return r.db.WithContext(ctx).Create(store).Error
}

func (r *storeRepositoryImpl) FindByUserID(ctx context.Context, userID string) (*entity.Store, error) {
	var store entity.Store
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&store).Error
	if err != nil {
		return nil, err
	}
	return &store, nil
}

func (r *storeRepositoryImpl) FindByID(ctx context.Context, id string) (*entity.Store, error) {
	var store entity.Store
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&store).Error
	if err != nil {
		return nil, err
	}
	return &store, nil
}

func (r *storeRepositoryImpl) Update(ctx context.Context, store *entity.Store) error {
	return r.db.WithContext(ctx).Save(store).Error
}
