package main

import "gorm.io/gorm"

type WalletRepository struct {
	DB *gorm.DB
}

func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{DB: db}
}

func (r *WalletRepository) CreateWallet(wallet *Wallet) error {
	return r.DB.Create(wallet).Error
}

func (r *WalletRepository) GetWallets() ([]Wallet, error) {
	var wallets []Wallet
	if err := r.DB.Find(&wallets).Error; err != nil {
		return nil, err
	}
	return wallets, nil
}

func (r *WalletRepository) DeleteWallet(address string) error {
	return r.DB.Where("address = ?", address).Delete(&Wallet{}).Error
}
