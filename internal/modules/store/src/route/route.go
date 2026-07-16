package route

import (
	"github.com/PetaFlops-web/backend-shop-smbk/internal/modules/store/src/controller"
	"github.com/gofiber/fiber/v2"
)

type StoreRouteConfig struct {
	App             *fiber.App
	StoreController *controller.StoreController
	JwtMiddleware   fiber.Handler
}

func (c *StoreRouteConfig) Setup() {
	api := c.App.Group("/api")
	
	stores := api.Group("/stores", c.JwtMiddleware)
	stores.Post("/", c.StoreController.Create)
	stores.Get("/", c.StoreController.GetMyStore)
	stores.Put("/", c.StoreController.UpdateMyStore)
}
