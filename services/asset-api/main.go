package main

import (
	"asset-management/pkg/app"
	"asset-management/pkg/database"
	"asset-management/pkg/logger"
	deposit2 "asset-management/services/asset-api/deposit"
	_ "asset-management/services/asset-api/docs"
	"asset-management/services/asset-api/transaction"
	"asset-management/services/asset-api/transfer"
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

const createBalanceTableSQL = `
CREATE TABLE IF NOT EXISTS balance (
	wallet_address VARCHAR(255) NOT NULL,
	network VARCHAR(100) NOT NULL,
	balance NUMERIC(30, 10) NOT NULL DEFAULT 0,
	UNIQUE (wallet_address, network)
);
`

const createScheduledTransactionsTable = `
CREATE TABLE IF NOT EXISTS scheduled_transactions (
    scheduled_transaction_id SERIAL PRIMARY KEY,
    from_wallet_address VARCHAR(255) NOT NULL,
    to_wallet_address VARCHAR(255) NOT NULL,
    network VARCHAR(100) NOT NULL,
    amount NUMERIC(30, 10) NOT NULL CHECK (amount > 0),
    scheduled_time TIMESTAMP NOT NULL,
    status VARCHAR(50) DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'COMPLETED', 'FAILED')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
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

	if err := CreateTables(db.Conn); err != nil {
		log.Fatal().Msgf("Failed to create balance table:", err)
	}

	appInstance := app.NewApp()

	appInstance.Fiber.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/swagger/index.html")
	})

	appInstance.AddRoute("/swagger/*", fiberSwagger.WrapHandler)

	walletValidator := wallet.NewValidationAdapter(os.Getenv("WALLET_API"))
	depositR := deposit2.NewDepositRepository(db.Conn)
	depositS := deposit2.NewDepositService(walletValidator, depositR)
	depositC := deposit2.NewDepositController(depositS)

	withdrawR := withdraw.NewWithdrawRepository(db.Conn)
	withdrawS := withdraw.NewWithdrawService(withdrawR, walletValidator)
	withdrawC := withdraw.NewWithdrawController(withdrawS)

	transferR := transfer.NewTransferRepository(db.Conn)
	transferS := transfer.NewTransferService(transferR, walletValidator)
	transferC := transfer.NewTransferController(transferS)

	scheduleTransferR := transaction.NewCreateRepository(db.Conn)
	scheduleTransferS := transaction.NewCreateService(scheduleTransferR, walletValidator)
	scheduleTransferC := transaction.NewCreateController(scheduleTransferS)

	appInstance.Fiber.Post("/deposit", depositC.Deposit)
	appInstance.Fiber.Post("/withdraw", withdrawC.Withdraw)
	appInstance.Fiber.Post("/transfer", transferC.Transfer)
	appInstance.Fiber.Post("/scheduled-transaction", scheduleTransferC.Create)
	//appInstance.Fiber.Get("/scheduled-transaction/next-minute", scheduleTransferC.GetNextMinuteTransactions)
	//appInstance.Fiber.Post("/scheduled-transaction/:id/process", scheduleTransferC.Process)

	log.Info().Msg("Asset Service is running on port 8081")
	appInstance.Start(":8001")
}

func CreateTables(db *sql.DB) error {
	// Execute the query
	_, err := db.Exec(createBalanceTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create balance table: %w", err)
	}

	_, schErr := db.Exec(createScheduledTransactionsTable)
	if schErr != nil {
		return fmt.Errorf("failed to create scheduled table: %w", err)
	}

	return nil
}
