package schedule

import (
	"context"
	"database/sql"
	"fmt"
)

type ScheduleTransactionRepository interface {
	Create(tx *ScheduleTransaction) (int, error)
	GetNextMinuteTransactions() ([]ScheduleTransaction, error)
	Process(scheduledTransactionID int64) error
}

type postgresScheduleTransactionRepository struct {
	db *sql.DB
}

func NewScheduleTransactionRepository(db *sql.DB) ScheduleTransactionRepository {
	return &postgresScheduleTransactionRepository{db: db}
}

// Create inserts a new schedule transaction into the database.
func (r *postgresScheduleTransactionRepository) Create(tx *ScheduleTransaction) (int, error) {
	query := `
		INSERT INTO scheduled_transactions (from_wallet_address, to_wallet_address, network, amount, scheduled_time, status)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING scheduled_transaction_id
	`
	var id int
	err := r.db.QueryRow(query, tx.FromWallet, tx.ToWallet, tx.Network, tx.Amount, tx.ScheduledTime, tx.Status).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to insert schedule transaction: %v", err)
	}
	return id, nil
}

func (r *postgresScheduleTransactionRepository) GetNextMinuteTransactions() ([]ScheduleTransaction, error) {
	rows, err := r.db.Query(`
        SELECT scheduled_transaction_id, from_wallet_address, to_wallet_address, network, amount, scheduled_time, status, created_at
        FROM scheduled_transactions
        WHERE scheduled_time >= NOW() AT TIME ZONE 'Europe/Istanbul'
          AND scheduled_time < NOW() AT TIME ZONE 'Europe/Istanbul' + INTERVAL '1 minute'
          AND status = 'PENDING'`)

	if err != nil {
		return nil, err
	}

	var transactions []ScheduleTransaction
	for rows.Next() {
		var txn ScheduleTransaction
		if err := rows.Scan(&txn.ID, &txn.FromWallet, &txn.ToWallet, &txn.Network, &txn.Amount, &txn.ScheduledTime, &txn.Status, &txn.CreatedAt); err != nil {
			return nil, err
		}
		transactions = append(transactions, txn)
	}

	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("failed to close rows: %w", err)
	}

	return transactions, nil
}

func (r *postgresScheduleTransactionRepository) Process(scheduledTransactionID int64) error {
	ctx := context.Background()

	// Begin transaction
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	// Ensure rollback in case of error
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
		return fmt.Errorf("insufficient funds in sender's wallet")
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
