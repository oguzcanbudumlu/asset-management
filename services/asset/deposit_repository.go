package main

import "database/sql"

type DepositRepository interface {
	Deposit() error
}

type depositRepository struct {
	db *sql.DB
}

func NewDepositRepository(db *sql.DB) DepositRepository {
	return &depositRepository{db: db}
}

func (*depositRepository) Deposit() error {
	return nil
}
