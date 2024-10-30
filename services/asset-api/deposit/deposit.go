package deposit

import (
	"github.com/shopspring/decimal"
)

// DepositRequest represents the request payload for a deposit
type DepositRequest struct {
	WalletAddress string          `json:"wallet_address" example:"0x123abc456def"`
	Network       string          `json:"network" example:"Ethereum"`
	Amount        decimal.Decimal `json:"amount" example:"100.50"`
}

// DepositResponse represents the response payload after a successful deposit
type DepositResponse struct {
	NewBalance decimal.Decimal `json:"new_balance" example:"1500.75"`
}
