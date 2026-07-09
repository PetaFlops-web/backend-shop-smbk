package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

func AuthMiddleware(viper *viper.Viper) fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		authHeader := ctx.Get("Authorization")
		if authHeader == "" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"errors": "Unauthorized",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"errors": "Unauthorized",
			})
		}

		tokenString := parts[1]
		jwtSecret := viper.GetString("jwt.secret")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"errors": "Unauthorized",
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"errors": "Unauthorized",
			})
		}

		userID := claims["id"].(string)

		ctx.Locals("auth", &Auth{ID: userID})

		return ctx.Next()
	}
}

func GetUser(ctx *fiber.Ctx) *Auth {
	return ctx.Locals("auth").(*Auth)
}
