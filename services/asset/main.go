package main

import (
	"asset-management/pkg/app"
	"asset-management/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// @title Asset Service API
// @version 1.0
// @description API documentation for Asset Service.
// @host localhost:8000
// @BasePath /
func main() {
	logger.InitLogger(zerolog.InfoLevel)

	appInstance := app.NewApp()

	appInstance.Fiber.Get("/asset", SimpleGet)

	appInstance.AddRoute("/swagger/*", fiberSwagger.WrapHandler)
	log.Info().Msg("Asset Service is running on port 8000")
	appInstance.Start(":8000")
}

// SimpleGet
// @Summary Get asset status
// @Description Returns the status of the asset service
// @Tags Asset
// @Accept json
// @Produce json
// @Success 200 {string} string "Asset Service is running"
// @Router /asset [get]
func SimpleGet(c *fiber.Ctx) error {
	log.Info().Msg("Handling asset request")
	return c.SendString("Asset Service is running")
}
