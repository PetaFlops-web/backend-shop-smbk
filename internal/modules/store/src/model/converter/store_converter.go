package converter

import (
	"github.com/PetaFlops-web/backend-shop-smbk/internal/modules/store/src/entity"
	"github.com/PetaFlops-web/backend-shop-smbk/internal/modules/store/src/model"
)

func StoreToResponse(store *entity.Store) *model.StoreResponse {
	if store == nil {
		return nil
	}
	return &model.StoreResponse{
		ID:        store.ID,
		UserID:    store.UserID,
		StoreName: store.StoreName,
		CreatedAt: store.CreatedAt,
		UpdatedAt: store.UpdatedAt,
	}
}