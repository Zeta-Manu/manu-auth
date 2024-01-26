package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"github.com/Zeta-Manu/manu-auth/internal/config"
	"github.com/Zeta-Manu/manu-auth/internal/infrastructure/healthcheck"
	"github.com/Zeta-Manu/manu-auth/internal/routes"
	"github.com/Zeta-Manu/manu-auth/internal/service"
)

func main() {
	appConfig, err := config.InitConfig("internal/config/config.json")
	if err != nil {
		fmt.Println("Error initializing configuration: ", err)
		return
	}

	app := fiber.New()

	userService, err := service.NewUserService(appConfig)
	if err != nil {
		fmt.Println("Error creating UserService: ", err)
		return
	}

	// Register routes
	healthcheck.RegisterRoutes(app)
	routes.RegisterUserRoutes(app, userService)

	app.Listen(":3000")
}
