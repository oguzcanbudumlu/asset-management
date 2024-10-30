package schedule

import (
	"database/sql"
	"fmt"
)

type ScheduleTransactionRepository interface {
	Create(tx *ScheduleTransaction) (int, error)
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
