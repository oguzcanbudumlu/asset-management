package main

import "database/sql"

type TransferRepository interface {
	Deposit() error
}

type transferRepository struct {
	db *sql.DB
}

func NewTransferRepository(db *sql.DB) TransferRepository {
	return &transferRepository{db: db}
}

func (*transferRepository) Deposit() error {
	return nil
}
