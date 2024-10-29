package main

import (
	"asset-management/pkg/app"
	"asset-management/pkg/database"
	"asset-management/pkg/logger"
	_ "asset-management/services/asset/docs"
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	"os"
)

const createBalanceTableSQL = `
CREATE TABLE IF NOT EXISTS balance (
	wallet_address VARCHAR(255) NOT NULL,
	network VARCHAR(100) NOT NULL,
	balance NUMERIC(18, 2) NOT NULL DEFAULT 0,
	UNIQUE (wallet_address, network)
);
`

// @title Asset Service API
// @version 1.0
// @description API documentation for Asset Service.
// @BasePath /
func main() {
	logger.InitLogger(zerolog.InfoLevel)
	db, err := database.NewDatabaseRaw(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))

	if err != nil {
		log.Error().Err(err).Msg("Failed to initialize database")
		return
	}
	defer db.Close()

	if err := CreateBalanceTable(db.Conn); err != nil {
		log.Fatal().Msgf("Failed to create balance table:", err)
	}

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

func CreateBalanceTable(db *sql.DB) error {
	// Execute the query
	_, err := db.Exec(createBalanceTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create balance table: %w", err)
	}

	return nil
}
