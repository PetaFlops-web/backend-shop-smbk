package auth

import (
	auth_client "github.com/PetaFlops-web/backend-shop-smbk/internal/module/auth-client"
	"github.com/PetaFlops-web/backend-shop-smbk/internal/module/auth/src/controller"
	"github.com/PetaFlops-web/backend-shop-smbk/internal/module/auth/src/entity"
	"github.com/PetaFlops-web/backend-shop-smbk/internal/module/auth/src/repository"
	"github.com/PetaFlops-web/backend-shop-smbk/internal/module/auth/src/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type Module struct {
	Controller *controller.AuthController
	UseCase    *usecase.AuthUseCase
	client     *clientImpl
	db         *gorm.DB
}

func New(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, config *viper.Viper) *Module {
	userRepo := repository.NewUserRepository(log)
	authUseCase := usecase.NewAuthUseCase(db, log, validate, userRepo, config)
	authController := controller.NewAuthController(authUseCase, log)

	return &Module{
		Controller: authController,
		UseCase:    authUseCase,
		client:     &clientImpl{db: db},
		db:         db,
	}
}

func (m *Module) Client() auth_client.Client {
	return m.client
}

func (m *Module) Migrate() error {
	return m.db.AutoMigrate(&entity.User{})
}