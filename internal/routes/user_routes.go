package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/Zeta-Manu/manu-auth/internal/domain"
	"github.com/Zeta-Manu/manu-auth/internal/service"
)

// RegisterUserRoutes registers user-related routes
func RegisterUserRoutes(app *fiber.App, userService *service.UserService) {
	v1 := app.Group("v1")

	// User registration route
	v1.Post("/register", func(c *fiber.Ctx) error {
		var userRegistration domain.UserRegistration
		if err := c.BodyParser(&userRegistration); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}

		if err := userService.RegisterUserHandler(userRegistration); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error registering user"})
		}

		return c.JSON(fiber.Map{"message": "User registered successfully"})
	})

	// User login route
	v1.Post("/login", func(c *fiber.Ctx) error {
		var userLogin domain.UserLogin
		if err := c.BodyParser(&userLogin); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}

		tokens, err := userService.LoginUserHandler(userLogin)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
		}

		return c.JSON(fiber.Map{"access_token": tokens.AccessToken, "id_token": tokens.IdToken, "refresh_token": tokens.RefreshToken})
	})
}
