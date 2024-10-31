package scheduled_process

import (
	"asset-management/internal/schedule"
	"context"
	"database/sql"
	"fmt"
	"github.com/rs/zerolog/log"
)

type ProcessRepository interface {
	Process(scheduledTransactionID int) error
}

type postgresProcessRepository struct {
	db *sql.DB
}

func NewProcessRepository(db *sql.DB) ProcessRepository {
	return &postgresProcessRepository{db: db}
}

func (r *postgresProcessRepository) Process(scheduledTransactionID int) error {
	ctx := context.Background()

	// Check if the transaction is already completed
	var status string
	err := r.db.QueryRowContext(ctx, `
        SELECT status FROM scheduled_transactions 
        WHERE scheduled_transaction_id = $1`, scheduledTransactionID).
		Scan(&status)
	if err != nil {
		return fmt.Errorf("failed to check transaction status: %v", err)
	}

	if status == schedule.StatusCompleted {
		// Early exit if transaction is already completed
		log.Info().Int("scheduledTransactionID", scheduledTransactionID).Msg("Transaction already completed, exiting early")
		return nil
	}

	// Begin transaction
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	// Define rollback function
	rollback := func() {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			fmt.Printf("transaction rollback failed: %v\n", rollbackErr)
		}
	}

	// Retrieve scheduled transaction details
	var fromWallet, toWallet, network string
	var amount float64

	err = tx.QueryRowContext(ctx, `
        SELECT from_wallet_address, to_wallet_address, network, amount 
        FROM scheduled_transactions 
        WHERE scheduled_transaction_id = $1 FOR UPDATE`, scheduledTransactionID).
		Scan(&fromWallet, &toWallet, &network, &amount)
	if err != nil {
		rollback()
		return fmt.Errorf("failed to fetch scheduled transaction: %v", err)
	}

	// Lock balance records for both from_wallet and to_wallet
	_, err = tx.ExecContext(ctx, `
        SELECT balance FROM balance 
        WHERE wallet_address = $1 AND network = $2 FOR UPDATE`, fromWallet, network)
	if err != nil {
		rollback()
		return fmt.Errorf("failed to lock sender's balance: %v", err)
	}

	_, err = tx.ExecContext(ctx, `
        SELECT balance FROM balance 
        WHERE wallet_address = $1 AND network = $2 FOR UPDATE`, toWallet, network)
	if err != nil {
		rollback()
		return fmt.Errorf("failed to lock receiver's balance: %v", err)
	}

	// Deduct from sender's balance
	res, err := tx.ExecContext(ctx, `
        UPDATE balance SET balance = balance - $1 
        WHERE wallet_address = $2 AND network = $3 AND balance >= $1`, amount, fromWallet, network)
	if err != nil {
		rollback()
		return fmt.Errorf("failed to deduct from sender's balance: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		rollback()
		return fmt.Errorf("insufficient balance in sender's wallet")
	}

	// Add to receiver's balance
	_, err = tx.ExecContext(ctx, `
    INSERT INTO balance (wallet_address, network, balance) 
    VALUES ($1, $2, $3) 
    ON CONFLICT (wallet_address, network) DO UPDATE 
    SET balance = balance.balance + EXCLUDED.balance`, toWallet, network, amount)

	if err != nil {
		rollback()
		return fmt.Errorf("failed to add to receiver's balance: %v", err)
	}

	// Update the scheduled transaction status to COMPLETED
	_, err = tx.ExecContext(ctx, `
        UPDATE scheduled_transactions SET status = 'COMPLETED' 
        WHERE scheduled_transaction_id = $1`, scheduledTransactionID)
	if err != nil {
		rollback()
		return fmt.Errorf("failed to update scheduled transaction status: %v", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		rollback()
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}
