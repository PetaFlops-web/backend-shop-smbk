package product

import (
	product_client "github.com/PetaFlops-web/backend-shop-smbk/internal/module/product-client"
	"github.com/PetaFlops-web/backend-shop-smbk/internal/module/product/src/controller"
	"github.com/PetaFlops-web/backend-shop-smbk/internal/module/product/src/entity"
	"github.com/PetaFlops-web/backend-shop-smbk/internal/module/product/src/repository"
	"github.com/PetaFlops-web/backend-shop-smbk/internal/module/product/src/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type Module struct {
	Controller *controller.ProductController
	client     *clientImpl
	db         *gorm.DB
}

func New(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, config *viper.Viper) *Module {
	productRepo := repository.NewProductRepository(log)
	productUseCase := usecase.NewProductUseCase(db, log, validate, productRepo)
	productController := controller.NewProductController(productUseCase, log)

	return &Module{
		Controller: productController,
		client:     &clientImpl{db: db, productRepo: productRepo},
		db:         db,
	}
}

func (m *Module) Client() product_client.Client {
	return m.client
}

func (m *Module) Migrate() error {
	return m.db.AutoMigrate(&entity.Product{})
}