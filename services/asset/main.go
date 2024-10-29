package main

import (
	"asset-management/pkg/app"
	"asset-management/pkg/database"
	"asset-management/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// @title Asset Service API
// @version 1.0
// @description API documentation for Asset Service.
// @host localhost:8081
// @BasePath /
func main() {
	logger.InitLogger(zerolog.InfoLevel)
	db, err := database.NewDatabaseRaw("localhost", "5431", "asset", "asset", "asset")
	if err != nil {
		log.Error().Err(err).Msg("Failed to initialize database")
		return
	}
	defer db.Close()

	appInstance := app.NewApp()

	appInstance.Fiber.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/swagger/index.html")
	})

	appInstance.AddRoute("/swagger/*", fiberSwagger.WrapHandler)

	walletValidator := NewWalletValidationAdapter("localhost:8080")
	depositR := NewDepositRepository(db.Conn)
	depositS := NewDepositService(walletValidator, depositR)
	depositC := NewDepositController(depositS)
	appInstance.Fiber.Post("/deposit", depositC.Deposit)
	log.Info().Msg("Asset Service is running on port 8081")
	appInstance.Start(":8001")
}
