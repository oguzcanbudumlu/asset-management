package withdraw

import (
	"database/sql"
	"errors"
	"github.com/shopspring/decimal"
)

type WithdrawRepository interface {
	Withdraw(walletAddress, network string, amount decimal.Decimal) (decimal.Decimal, error)
}

type withdrawRepository struct {
	db *sql.DB
}

func NewWithdrawRepository(db *sql.DB) WithdrawRepository {
	return &withdrawRepository{db: db}
}

func (r *withdrawRepository) Withdraw(walletAddress, network string, amount decimal.Decimal) (decimal.Decimal, error) {
	var currentBalance decimal.Decimal

	// Start a transaction
	tx, err := r.db.Begin()
	if err != nil {
		return decimal.Zero, err
	}
	defer tx.Rollback()

	// Check the current balance
	err = tx.QueryRow(`
        SELECT balance 
        FROM balance 
        WHERE wallet_address = $1 AND network = $2 
        FOR UPDATE`, walletAddress, network).Scan(&currentBalance)

	if err == sql.ErrNoRows {
		return decimal.Zero, errors.New("wallet not found")
	} else if err != nil {
		return decimal.Zero, err
	}

	// Check if balance is sufficient
	if currentBalance.LessThan(amount) {
		return decimal.Zero, errors.New("insufficient balance")
	}

	// Perform the withdrawal by updating the balance
	_, err = tx.Exec(`
        UPDATE balance 
        SET balance = balance - $1 
        WHERE wallet_address = $2 AND network = $3`, amount, walletAddress, network)

	if err != nil {
		return decimal.Zero, err
	}

	// Commit the transaction and get the new balance
	err = tx.Commit()
	if err != nil {
		return decimal.Zero, err
	}

	// Return the updated balance
	return currentBalance.Sub(amount), nil
}
