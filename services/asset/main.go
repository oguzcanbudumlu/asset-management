package main

import (
	"asset-management/pkg/app"
	"asset-management/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	logger.InitLogger(zerolog.InfoLevel)

	appInstance := app.NewApp()

	appInstance.AddRoute("/asset", func(c *fiber.Ctx) error {
		log.Info().Msg("Handling asset request")
		return c.SendString("Asset Service is running")
	})

	log.Info().Msg("Asset Service is running on port 8000")
	appInstance.Start(":8000")
}
