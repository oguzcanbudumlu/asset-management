package withdraw

import "database/sql"

type WithdrawRepository interface {
	Deposit() error
}

type withdrawRepository struct {
	db *sql.DB
}

func NewWithdrawRepository(db *sql.DB) WithdrawRepository {
	return &withdrawRepository{db: db}
}

func (*withdrawRepository) Deposit() error {
	return nil
}
