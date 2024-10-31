package transfer

import (
	"database/sql"
	"errors"
	"fmt"
)

type TransferRepository interface {
	Transfer(from, to, network string, amount float64) error
}

type transferRepository struct {
	db *sql.DB
}

func NewTransferRepository(db *sql.DB) TransferRepository {
	return &transferRepository{db: db}
}

func (r *transferRepository) Transfer(fromWallet, toWallet, network string, amount float64) error {
	// Start transaction
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	// Lock and check the sender's current balance
	var fromBalance float64
	query := `SELECT balance FROM balance WHERE wallet_address = $1 AND network = $2 FOR UPDATE`
	if err = tx.QueryRow(query, fromWallet, network).Scan(&fromBalance); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("failed to fetch balance for fromWallet: %v, rollback error: %v", err, rbErr)
		}
		return fmt.Errorf("failed to fetch balance for fromWallet: %v", err)
	}

	// Check if there is sufficient balance
	if fromBalance < amount {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("insufficient balance in source wallet: %v, rollback error: %v", errors.New("insufficient balance"), rbErr)
		}
		return errors.New("insufficient balance in source wallet")
	}

	// Deduct balance from the sender's wallet
	updateFrom := `UPDATE balance SET balance = balance - $1 WHERE wallet_address = $2 AND network = $3`
	if _, err = tx.Exec(updateFrom, amount, fromWallet, network); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("failed to update balance for fromWallet: %v, rollback error: %v", err, rbErr)
		}
		return fmt.Errorf("failed to update balance for fromWallet: %v", err)
	}

	// Check if the recipient's balance exists (if not, add it)
	var toBalanceExists bool
	checkTo := `SELECT EXISTS (SELECT 1 FROM balance WHERE wallet_address = $1 AND network = $2)`
	if err = tx.QueryRow(checkTo, toWallet, network).Scan(&toBalanceExists); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("failed to check balance existence for toWallet: %v, rollback error: %v", err, rbErr)
		}
		return fmt.Errorf("failed to check balance existence for toWallet: %v", err)
	}

	// Update recipient's balance if it exists; otherwise, insert a new record
	if toBalanceExists {
		updateTo := `UPDATE balance SET balance = balance + $1 WHERE wallet_address = $2 AND network = $3`
		if _, err = tx.Exec(updateTo, amount, toWallet, network); err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				return fmt.Errorf("failed to update balance for toWallet: %v, rollback error: %v", err, rbErr)
			}
			return fmt.Errorf("failed to update balance for toWallet: %v", err)
		}
	} else {
		insertTo := `INSERT INTO balance (wallet_address, network, balance) VALUES ($1, $2, $3)`
		if _, err = tx.Exec(insertTo, toWallet, network, amount); err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				return fmt.Errorf("failed to insert balance for toWallet: %v, rollback error: %v", err, rbErr)
			}
			return fmt.Errorf("failed to insert balance for toWallet: %v", err)
		}
	}

	// Commit the transaction if successful
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil // Transaction successful
}
