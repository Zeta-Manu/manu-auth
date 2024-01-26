package healthcheck

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app *fiber.App) {
	app.Get("/health", healthCheckHandler)
}

func healthCheckHandler(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusOK)
}
