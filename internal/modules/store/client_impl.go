package store

import (
	"github.com/PetaFlops-web/backend-shop-smbk/internal/modules/store/src/repository"
)

type clientImpl struct {
	repo repository.StoreRepository
}

