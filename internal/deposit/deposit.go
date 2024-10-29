package deposit

type DepositRequest struct {
	WalletAddress string  `json:"wallet_address" example:"0x123abc456def"`
	Network       string  `json:"network" example:"Ethereum"`
	Amount        float64 `json:"amount" example:"100.50"`
}

type DepositResponse struct {
	TransactionID string  `json:"transaction_id" example:"txn_123456789"`
	NewBalance    float64 `json:"new_balance" example:"1500.75"`
}
