package product

import (
	"github.com/gofiber/fiber/v2"
)

func (m *Module) RegisterRoutes(router fiber.Router, authMiddleware fiber.Handler) {
	productRouter := router.Group("/products")
	productRouter.Use(authMiddleware)

	productRouter.Post("/", m.Controller.Create)
	productRouter.Get("/", m.Controller.List)
	productRouter.Get("/:id", m.Controller.Get)
	productRouter.Put("/:id", m.Controller.Update)
	productRouter.Delete("/:id", m.Controller.Delete)
}