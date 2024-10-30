package transfer

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
)

type TransferRepository interface {
	Transfer(from, to, network string, amount decimal.Decimal) error
}

type transferRepository struct {
	db *sql.DB
}

func NewTransferRepository(db *sql.DB) TransferRepository {
	return &transferRepository{db: db}
}

func (r *transferRepository) Transfer(fromWallet, toWallet, network string, amount decimal.Decimal) error {
	// Start transaction
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	// Define inner rollback function
	rollback := func(action string, originalErr error) error {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("%s: %v, rollback error: %v", action, originalErr, rbErr)
		}
		return fmt.Errorf("%s: %v", action, originalErr)
	}

	// Lock and check the sender's current balance
	var fromBalance decimal.Decimal
	query := `SELECT balance FROM balance WHERE wallet_address = $1 AND network = $2 FOR UPDATE`
	if err = tx.QueryRow(query, fromWallet, network).Scan(&fromBalance); err != nil {
		return rollback("failed to fetch balance for fromWallet", err)
	}

	// Check if there is sufficient balance
	if fromBalance.LessThan(amount) {
		return rollback("insufficient balance in source wallet", errors.New("insufficient balance"))
	}

	// Deduct balance from the sender's wallet
	updateFrom := `UPDATE balance SET balance = balance - $1 WHERE wallet_address = $2 AND network = $3`
	if _, err = tx.Exec(updateFrom, amount, fromWallet, network); err != nil {
		return rollback("failed to update balance for fromWallet", err)
	}

	// Check if the recipient's balance exists
	var toBalanceExists bool
	checkTo := `SELECT EXISTS (SELECT 1 FROM balance WHERE wallet_address = $1 AND network = $2)`
	if err = tx.QueryRow(checkTo, toWallet, network).Scan(&toBalanceExists); err != nil {
		return rollback("failed to check balance existence for toWallet", err)
	}

	// Update recipient's balance if it exists; otherwise, insert a new record
	if toBalanceExists {
		updateTo := `UPDATE balance SET balance = balance + $1 WHERE wallet_address = $2 AND network = $3`
		if _, err = tx.Exec(updateTo, amount, toWallet, network); err != nil {
			return rollback("failed to update balance for toWallet", err)
		}
	} else {
		insertTo := `INSERT INTO balance (wallet_address, network, balance) VALUES ($1, $2, $3)`
		if _, err = tx.Exec(insertTo, toWallet, network, amount); err != nil {
			return rollback("failed to insert balance for toWallet", err)
		}
	}

	// Commit the transaction if successful
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil // Transaction successful
}
