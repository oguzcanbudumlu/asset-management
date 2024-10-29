package main

import (
	"database/sql"
	"fmt"
)

type DepositRepository interface {
	Deposit(walletAddress, network string, amount float64) (float64, error)
}

type depositRepository struct {
	db *sql.DB
}

func NewDepositRepository(db *sql.DB) DepositRepository {
	return &depositRepository{db: db}
}

func (r *depositRepository) Deposit(walletAddress, network string, amount float64) (float64, error) {
	// Start a new transaction
	tx, err := r.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// Use UPSERT to insert or update the balance atomically
	upsertQuery := `
		INSERT INTO balance (wallet_address, network, balance)
		VALUES ($1, $2, $3)
		ON CONFLICT (wallet_address, network) 
		DO UPDATE SET balance = balance + EXCLUDED.balance
		RETURNING balance;
	`

	var newBalance float64
	err = tx.QueryRow(upsertQuery, walletAddress, network, amount).Scan(&newBalance)
	if err != nil {
		return 0, fmt.Errorf("failed to upsert balance: %w", err)
	}

	return newBalance, nil
}
