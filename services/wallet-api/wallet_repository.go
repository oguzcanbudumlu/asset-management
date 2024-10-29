package main

import (
	"fmt"
	"gorm.io/gorm"
)

// WalletRepository defines the methods for wallet data persistence
type WalletRepository interface {
	CreateWallet(wallet *Wallet) error
	GetWallets() ([]Wallet, error)
	DeleteWallet(network, address string) error
	GetWallet(network, address string) (*Wallet, error)
}

type walletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) WalletRepository {
	return &walletRepository{db: db}
}

func (r *walletRepository) CreateWallet(wallet *Wallet) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Check if the wallet with the same address and network already exists
		var count int64
		if err := tx.Model(&Wallet{}).Where("address = ? AND network = ?", wallet.Address, wallet.Network).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return fmt.Errorf("wallet with address %s and network %s already exists", wallet.Address, wallet.Network)
		}

		// If unique, create the wallet
		return tx.Create(wallet).Error
	})
}

func (r *walletRepository) GetWallets() ([]Wallet, error) {
	var wallets []Wallet
	if err := r.db.Find(&wallets).Error; err != nil {
		return nil, err
	}
	return wallets, nil
}

func (r *walletRepository) DeleteWallet(network, address string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Move the wallet to wallet_deleted table
		var wallet Wallet
		if err := tx.Where("network = ? AND address = ?", network, address).First(&wallet).Error; err != nil {
			return err // Return error if the wallet is not found
		}

		// Insert into wallet_deleted
		deletedWallet := WalletDeleted{
			ID:      wallet.ID,
			Network: wallet.Network,
			Address: wallet.Address,
		}
		if err := tx.Create(&deletedWallet).Error; err != nil {
			return err // Return error if unable to insert into wallet_deleted
		}

		// Now delete the original wallet
		return tx.Delete(&wallet).Error
	})
}

func (r *walletRepository) GetWallet(network, address string) (*Wallet, error) {
	var wallet Wallet
	if err := r.db.Where("network = ? AND address = ?", network, address).First(&wallet).Error; err != nil {
		return nil, err
	}
	return &wallet, nil
}
