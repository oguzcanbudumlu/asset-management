package scheduled_next

import (
	"asset-management/internal/schedule"
	"database/sql"
	"fmt"
)

type NextRepository interface {
	GetNextMinuteTransactions() ([]schedule.ScheduledTransaction, error)
}

type postgresNextRepository struct {
	db *sql.DB
}

func NewNextRepository(db *sql.DB) NextRepository {
	return &postgresNextRepository{db: db}
}

func (r *postgresNextRepository) GetNextMinuteTransactions() ([]schedule.ScheduledTransaction, error) {
	rows, err := r.db.Query(`
        SELECT scheduled_transaction_id, from_wallet_address, to_wallet_address, network, amount, scheduled_time, status, created_at
        FROM scheduled_transactions
        WHERE scheduled_time >= (NOW() AT TIME ZONE 'Europe/Istanbul' - INTERVAL '5 minute')
          AND scheduled_time < (NOW() AT TIME ZONE 'Europe/Istanbul' + INTERVAL '5 minute')
          AND status = 'PENDING'`)

	if err != nil {
		return nil, err
	}

	var transactions []schedule.ScheduledTransaction
	for rows.Next() {
		var txn schedule.ScheduledTransaction
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
