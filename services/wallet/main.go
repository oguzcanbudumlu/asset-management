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

	appInstance.AddRoute("/wallet", func(c *fiber.Ctx) error {
		log.Info().Msg("Handling wallet request")
		return c.SendString("Wallet Service is running")
	})

	log.Info().Msg("Wallet Service is running on port 8001")
	appInstance.Start(":8001")
}
