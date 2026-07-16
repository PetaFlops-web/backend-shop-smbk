package usecase

import (
	"context"
	"net/http"

	"github.com/PetaFlops-web/backend-shop-smbk/internal/modules/store/src/entity"
	"github.com/PetaFlops-web/backend-shop-smbk/internal/modules/store/src/model"
	"github.com/PetaFlops-web/backend-shop-smbk/internal/modules/store/src/repository"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StoreUseCase interface {
	Create(ctx context.Context, userID string, request *model.CreateStoreRequest) (*model.StoreResponse, error)
	GetMyStore(ctx context.Context, userID string) (*model.StoreResponse, error)
	UpdateMyStore(ctx context.Context, userID string, request *model.UpdateStoreRequest) (*model.StoreResponse, error)
}

type storeUseCaseImpl struct {
	repo     repository.StoreRepository
	validate *validator.Validate
}

func NewStoreUseCase(repo repository.StoreRepository, validate *validator.Validate) StoreUseCase {
	return &storeUseCaseImpl{
		repo:     repo,
		validate: validate,
	}
}

func (u *storeUseCaseImpl) Create(ctx context.Context, userID string, request *model.CreateStoreRequest) (*model.StoreResponse, error) {
	if err := u.validate.Struct(request); err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	existingStore, err := u.repo.FindByUserID(ctx, userID)
	if err == nil && existingStore != nil {
		return nil, fiber.NewError(fiber.StatusConflict, "User already has a store")
	} else if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	store := &entity.Store{
		ID:        uuid.New().String(),
		UserID:    userID,
		StoreName: request.StoreName,
	}

	if err := u.repo.Create(ctx, store); err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to create store")
	}

	return toStoreResponse(store), nil
}

func (u *storeUseCaseImpl) GetMyStore(ctx context.Context, userID string) (*model.StoreResponse, error) {
	store, err := u.repo.FindByUserID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fiber.NewError(fiber.StatusNotFound, "Store not found")
		}
		return nil, err
	}

	return toStoreResponse(store), nil
}

func (u *storeUseCaseImpl) UpdateMyStore(ctx context.Context, userID string, request *model.UpdateStoreRequest) (*model.StoreResponse, error) {
	if err := u.validate.Struct(request); err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	store, err := u.repo.FindByUserID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fiber.NewError(fiber.StatusNotFound, "Store not found")
		}
		return nil, err
	}

	store.StoreName = request.StoreName

	if err := u.repo.Update(ctx, store); err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to update store")
	}

	return toStoreResponse(store), nil
}

func toStoreResponse(store *entity.Store) *model.StoreResponse {
	return &model.StoreResponse{
		ID:        store.ID,
		UserID:    store.UserID,
		StoreName: store.StoreName,
		CreatedAt: store.CreatedAt,
		UpdatedAt: store.UpdatedAt,
	}
}
