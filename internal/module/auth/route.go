package auth

import (
	"github.com/gofiber/fiber/v2"
)

func (m *Module) RegisterRoutes(router fiber.Router, authMiddleware fiber.Handler) {
	authRouter := router.Group("/auth")
	authRouter.Post("/register", m.Controller.Register)
	authRouter.Post("/login", m.Controller.Login)

	authRouter.Use(authMiddleware)
	authRouter.Get("/current", m.Controller.Current)
}