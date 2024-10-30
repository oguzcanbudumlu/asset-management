package main

import (
	"asset-management/pkg/app"
	"asset-management/pkg/database"
	"asset-management/pkg/logger"
	_ "asset-management/services/wallet-api/docs"
	wallet2 "asset-management/services/wallet-api/wallet"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	"os"
)

// @title Wallet Service API
// @version 1.0
// @description API documentation for Wallet Service.
// @BasePath /
func main() {
	logger.InitLogger(zerolog.DebugLevel)
	db, err := database.NewDatabase(os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to initialize database")
		return
	}

	appInstance := app.NewApp()

	appInstance.Fiber.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/swagger/index.html")
	})

	appInstance.Fiber.Get("/swagger/*", fiberSwagger.WrapHandler)
	if err := db.Conn.AutoMigrate(&wallet2.Wallet{}, &wallet2.WalletDeleted{}); err != nil {
		log.Error().Err(err).Msg("Failed to migrate database schema")
		return
	}

	repo := wallet2.NewWalletRepository(db.Conn)
	service := wallet2.NewWalletService(repo)
	controller := wallet2.NewWalletController(service)

	appInstance.Fiber.Post("/wallet", controller.CreateWallet)
	appInstance.Fiber.Get("/wallet/:network/:address", controller.GetWallet)
	appInstance.Fiber.Delete("/wallet/:network/:address", controller.DeleteWallet)

	log.Info().Msg("Wallet Service is running on port 8000")
	appInstance.Start(":8000")

	if dbCloseErr := db.Close(); dbCloseErr != nil {
		log.Error().Err(dbCloseErr).Msg("Failed to close database connection")
	}
}
