package usecase

import (
	"context"

	"github.com/PetaFlops-web/backend-shop-smbk/internal/module/auth/src/entity"
	"github.com/PetaFlops-web/backend-shop-smbk/internal/module/auth/src/model"
	"github.com/PetaFlops-web/backend-shop-smbk/internal/module/auth/src/model/converter"
	"github.com/PetaFlops-web/backend-shop-smbk/internal/module/auth/src/repository"
	"github.com/PetaFlops-web/backend-shop-smbk/internal/shared/middleware"
	"github.com/PetaFlops-web/backend-shop-smbk/internal/shared/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthUseCase struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	Validate       *validator.Validate
	UserRepository *repository.UserRepository
	Viper          *viper.Viper
}

func NewAuthUseCase(
	db *gorm.DB,
	log *logrus.Logger,
	validate *validator.Validate,
	userRepo *repository.UserRepository,
	viper *viper.Viper,
) *AuthUseCase {
	return &AuthUseCase{
		DB:             db,
		Log:            log,
		Validate:       validate,
		UserRepository: userRepo,
		Viper:          viper,
	}
}

func (c *AuthUseCase) Current(ctx context.Context, id string) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx)

	if err := c.Validate.Var(id, "required"); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	user := new(entity.User)
	if err := c.UserRepository.FindById(tx, user, id); err != nil {
		c.Log.Warnf("Failed find user by id : %+v", err)
		return nil, fiber.NewError(fiber.StatusNotFound, "User tidak ditemukan")
	}

	return converter.UserToResponse(user), nil
}

func (u *AuthUseCase) Register(ctx context.Context, request *model.RegisterRequest) (*model.AuthResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	total, err := u.UserRepository.CountByUsername(tx, request.Username)
	if err != nil {
		u.Log.Warnf("Failed count user from database : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to register user")
	}

	if total > 0 {
		u.Log.Warnf("Username already exists : %+v", request.Username)
		return nil, fiber.NewError(fiber.StatusConflict, "Username already taken")
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		u.Log.Warnf("Failed to generate bcrypt hash : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to register user")
	}

	userId, err := utils.GenerateUserId(request.Username)
	if err != nil {
		u.Log.Warnf("Failed to generate user id : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to register user")
	}

	user := &entity.User{
		ID:       userId,
		Password: string(password),
		Username: request.Username,
		Email:    request.Email,
	}

	if err := u.UserRepository.Create(tx, user); err != nil {
		u.Log.Warnf("Failed create user to database : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to register user")
	}

	token, err := middleware.GenerateJWT(u.Viper.GetString("jwt.secret"), user.ID, user.Username)
	if err != nil {
		u.Log.Errorf("Failed to generate JWT for user %s: %v", user.ID, err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to generate token")
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to register user")
	}

	return &model.AuthResponse{
		User:  *converter.UserToResponse(user),
		Token: token,
	}, nil
}

func (u *AuthUseCase) Login(ctx context.Context, request *model.LoginRequest) (*model.AuthResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	user := new(entity.User)
	if err := u.UserRepository.FindByUsername(tx, user, request.Username); err != nil {
		u.Log.Warnf("Failed find user by username : %+v", err)
		return nil, fiber.NewError(fiber.StatusNotFound, "Username atau password anda salah")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		u.Log.Warnf("Failed to compare user password with bcrypt hash : %+v", err)
		return nil, fiber.NewError(fiber.StatusNotFound, "Username atau password anda salah")
	}

	token, err := middleware.GenerateJWT(u.Viper.GetString("jwt.secret"), user.ID, user.Username)
	if err != nil {
		u.Log.Errorf("Failed to generate JWT for user %s: %v", user.ID, err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to generate token")
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to login")
	}

	return &model.AuthResponse{
		User:  *converter.UserToResponse(user),
		Token: token,
	}, nil
}