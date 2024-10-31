package withdraw

import (
	"database/sql"
	"errors"
)

type Repository interface {
	Withdraw(walletAddress, network string, amount float64) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Withdraw(walletAddress, network string, amount float64) error {
	var currentBalance float64

	// Start a transaction
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Check the current balance
	err = tx.QueryRow(`
        SELECT balance 
        FROM balance 
        WHERE wallet_address = $1 AND network = $2 
        FOR UPDATE`, walletAddress, network).Scan(&currentBalance)

	if err == sql.ErrNoRows {
		return errors.New("wallet not found")
	} else if err != nil {
		return err
	}

	// Check if balance is sufficient
	if currentBalance < amount {
		return errors.New("insufficient balance")
	}

	// Perform the withdrawal by updating the balance
	_, err = tx.Exec(`
        UPDATE balance 
        SET balance = balance - $1 
        WHERE wallet_address = $2 AND network = $3`, amount, walletAddress, network)

	if err != nil {
		return err
	}

	// Commit the transaction
	return tx.Commit()
}
