package main

import (
	"fmt"

	"github.com/PetaFlops-web/backend-shop-smbk/internal/module/auth"
	"github.com/PetaFlops-web/backend-shop-smbk/internal/shared/config"
	"github.com/PetaFlops-web/backend-shop-smbk/internal/shared/middleware"
	module "github.com/PetaFlops-web/backend-shop-smbk/internal/shared/module"
)

// @title           AIC Backend API
// @version         1.0
// @description     API Documentation for AIC Backend.

// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	db := config.NewDatabase(viperConfig, log)
	validate := config.NewValidator(viperConfig)
	app := config.NewFiber(viperConfig)


	// Auth Middleware
	authMiddleware := middleware.AuthMiddleware(viperConfig)

	// Module initialization (ordered by dependency)
	authModule := auth.New(db, log, validate, viperConfig)

	// Register all modules
	modules := []module.Module{
		authModule,
	}

	// Auto-migration (each module migrates its own tables)
	for _, m := range modules {
		if err := m.Migrate(); err != nil {
			log.Fatalf("Failed to migrate: %v", err)
		}
	}

	// Route registration
	api := app.Group("/api")
	for _, m := range modules {
		m.RegisterRoutes(api, authMiddleware)
	}

	// Start server
	port := viperConfig.GetInt("web.port")
	if port == 0 {
		port = 8080
	}

	log.Infof("Server is starting on port :%d", port)

	err := app.Listen(fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}