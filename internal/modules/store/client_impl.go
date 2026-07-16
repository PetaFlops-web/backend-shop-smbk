package store

import (
	"context"

	store_client "github.com/PetaFlops-web/backend-shop-smbk/internal/modules/store-client"
	"github.com/PetaFlops-web/backend-shop-smbk/internal/modules/store/src/repository"
)

type clientImpl struct {
	repo repository.StoreRepository
}

func (c *clientImpl) GetStoreByUserID(ctx context.Context, userID string) (*store_client.StoreResponse, error) {
	store, err := c.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &store_client.StoreResponse{
		ID:        store.ID,
		UserID:    store.UserID,
		StoreName: store.StoreName,
		CreatedAt: store.CreatedAt,
		UpdatedAt: store.UpdatedAt,
	}, nil
}

func (c *clientImpl) GetStoreByID(ctx context.Context, storeID string) (*store_client.StoreResponse, error) {
	store, err := c.repo.FindByID(ctx, storeID)
	if err != nil {
		return nil, err
	}

	return &store_client.StoreResponse{
		ID:        store.ID,
		UserID:    store.UserID,
		StoreName: store.StoreName,
		CreatedAt: store.CreatedAt,
		UpdatedAt: store.UpdatedAt,
	}, nil
}
