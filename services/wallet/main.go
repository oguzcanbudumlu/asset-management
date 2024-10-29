package main

import (
	"asset-management/pkg/app"
	"asset-management/pkg/database"
	"asset-management/pkg/logger"
	_ "asset-management/services/wallet/docs"
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
	db, err := newDb()
	if err != nil {
		log.Error().Err(err).Msg("Failed to initialize database")
		return
	}
	defer db.Close()

	appInstance := app.NewApp()

	appInstance.Fiber.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/swagger/index.html")
	})

	appInstance.Fiber.Get("/swagger/*", fiberSwagger.WrapHandler)
	if err := db.Conn.AutoMigrate(&Wallet{}, &WalletDeleted{}); err != nil {
		log.Error().Err(err).Msg("Failed to migrate database schema")
		return
	}

	repo := NewWalletRepository(db.Conn)
	service := NewWalletService(repo)
	controller := NewWalletController(service)

	appInstance.Fiber.Post("/wallet", controller.CreateWallet)
	appInstance.Fiber.Get("/wallet/:network/:address", controller.GetWallet)
	appInstance.Fiber.Delete("/wallet/:network/:address", controller.DeleteWallet)

	log.Info().Msg("Wallet Service is running on port 8000")
	appInstance.Start(":8000")
}

func newDb() (*database.Database, error) {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	return database.NewDatabase(dbHost, dbPort, dbUser, dbPassword, dbName)
}
