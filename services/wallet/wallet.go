package main

type Wallet struct {
	ID      uint   `gorm:"primaryKey"`
	Address string `gorm:"unique;not null"`
	Network string `gorm:"not null"`
}
