package schedule

import (
	"database/sql"
	"fmt"
)

type ScheduleTransactionRepository interface {
	Create(tx *ScheduleTransaction) (int, error)
	GetNextMinuteTransactions() ([]ScheduleTransaction, error)
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
          AND scheduled_time < NOW() AT TIME ZONE 'Europe/Istanbul' + INTERVAL '1 minute'`)

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
