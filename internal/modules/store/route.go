package store

import (
	"github.com/gofiber/fiber/v2"
)

func (m *Module) RegisterRoutes(router fiber.Router, authMiddleware fiber.Handler) {
	storeRouter := router.Group("/stores")
	storeRouter.Use(authMiddleware)

	storeRouter.Post("/", m.Controller.Create)
	storeRouter.Get("/", m.Controller.Get)
	storeRouter.Put("/", m.Controller.Update)
}
