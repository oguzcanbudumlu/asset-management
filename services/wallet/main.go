package main

import (
	"asset-management/pkg/app"
	"asset-management/pkg/logger"
	_ "asset-management/services/wallet/docs"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// @title Wallet Service API
// @version 1.0
// @description API documentation for Wallet Service.
// @host localhost:8001
// @BasePath /
func main() {
	logger.InitLogger(zerolog.InfoLevel)

	appInstance := app.NewApp()

	appInstance.Fiber.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/swagger/index.html")
	})
	appInstance.AddRoute("/wallet", simpleGet)

	appInstance.Fiber.Get("/swagger/*", fiberSwagger.WrapHandler)

	log.Info().Msg("Wallet Service is running on port 8001")
	appInstance.Start(":8001")
}

// simpleGet is an endpoint to get the status of the Wallet Service
// @Summary Get wallet status
// @Description Returns the status of the wallet service
// @Tags Wallet
// @Accept json
// @Produce json
// @Success 200 {string} string "Wallet Service is running"
// @Router /wallet [get]
func simpleGet(c *fiber.Ctx) error {
	log.Info().Msg("Handling wallet request")
	return c.SendString("Wallet Service is running")
}
