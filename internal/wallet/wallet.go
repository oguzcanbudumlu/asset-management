package wallet

// Wallet struct
// @Description Represents a wallet in the system
type Wallet struct {
	ID      uint   `gorm:"primaryKey;autoIncrement" json:"-"`
	Network string `gorm:"not null" json:"network"`
	Address string `gorm:"not null" json:"address"`
}

// WalletDeleted struct
// @Description Represents a deleted wallet in the system
type WalletDeleted struct {
	ID      uint   `gorm:"primaryKey" json:"-"`
	Network string `gorm:"not null" json:"network"`
	Address string `gorm:"not null" json:"address"`
}
