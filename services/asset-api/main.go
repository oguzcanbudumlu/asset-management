package main

import (
	sql2 "asset-management/internal/sql"
	"asset-management/pkg/app"
	"asset-management/pkg/database"
	"asset-management/pkg/logger"
	deposit2 "asset-management/services/asset-api/deposit"
	_ "asset-management/services/asset-api/docs"
	"asset-management/services/asset-api/scheduled"
	"asset-management/services/asset-api/wallet"
	"asset-management/services/asset-api/withdraw"
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	"os"
)

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

	if err := CreateTables(db.Conn); err != nil {
		log.Fatal().Err(err).Msg("Failed to create tables")
		dbCloseErr := db.Close()
		if dbCloseErr != nil {
			log.Error().Err(dbCloseErr).Msg("Failed to close database connection after error")
		}
		return
	}

	appInstance := app.NewApp()
	appInstance.Fiber.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/swagger/index.html")
	})

	appInstance.AddRoute("/swagger/*", fiberSwagger.WrapHandler)

	walletValidator := wallet.NewValidationAdapter(os.Getenv("WALLET_API"))
	depositR := deposit2.NewRepository(db.Conn)
	depositS := deposit2.NewService(walletValidator, depositR)
	depositC := deposit2.NewController(depositS)

	withdrawR := withdraw.NewRepository(db.Conn)
	withdrawS := withdraw.NewService(withdrawR, walletValidator)
	withdrawC := withdraw.NewController(withdrawS)

	scheduleTransferR := scheduled.NewCreateRepository(db.Conn)
	scheduleTransferS := scheduled.NewCreateService(scheduleTransferR, walletValidator)
	scheduleTransferC := scheduled.NewCreateController(scheduleTransferS)

	appInstance.Fiber.Post("/deposit", depositC.Deposit)
	appInstance.Fiber.Post("/withdraw", withdrawC.Withdraw)
	appInstance.Fiber.Post("/scheduled-transaction", scheduleTransferC.Create)
	//appInstance.Fiber.Get("/scheduled-transaction/next-minute", scheduleTransferC.GetNextMinuteTransactions)
	//appInstance.Fiber.Post("/scheduled-transaction/:id/process", scheduleTransferC.Process)

	log.Info().Msg("Asset Service is running on port 8081")
	appInstance.Start(":8001")

	if dbCloseErr := db.Close(); dbCloseErr != nil {
		log.Error().Err(dbCloseErr).Msg("Failed to close database connection")
	}

}

func CreateTables(db *sql.DB) error {
	if _, err := db.Exec(sql2.CreateBalanceTableSQL); err != nil {
		return fmt.Errorf("failed to create balance table: %w", err)
	}

	if _, schErr := db.Exec(sql2.CreateScheduledTransactionsTable); schErr != nil {
		return fmt.Errorf("failed to create scheduled transactions table: %w", schErr)
	}

	return nil
}
