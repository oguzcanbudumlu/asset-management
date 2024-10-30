package schedule

import "time"

type ScheduleTransaction struct {
	ID            int       `json:"id" example:"1"`                                // Transaction ID
	FromWallet    string    `json:"from_wallet" example:"wallet_123"`              // Sender's wallet address
	ToWallet      string    `json:"to_wallet" example:"wallet_456"`                // Recipient's wallet address
	Network       string    `json:"network" example:"Ethereum"`                    // Blockchain network (e.g., Ethereum)
	Amount        float64   `json:"amount" example:"250.75"`                       // Amount to be transferred
	ScheduledTime time.Time `json:"scheduled_time" example:"2024-10-30T15:04:05Z"` // Scheduled time for transaction
	Status        string    `json:"status" example:"PENDING"`                      // Transaction status (e.g., pending, completed)
	CreatedAt     time.Time `json:"created_at" example:"2024-10-29T10:15:00Z"`     // Time when the transaction was created
}

const (
	StatusPending   = "PENDING"
	StatusCompleted = "COMPLETED"
	StatusFailed    = "FAILED"
)

type Request struct {
	From          string  `json:"from" example:"wallet123"`
	To            string  `json:"to" example:"wallet456"`
	Network       string  `json:"network" example:"mainnet"`
	Amount        float64 `json:"amount" example:"100.50"`
	ScheduledTime string  `json:"scheduled_time" example:"2023-12-31T12:00:00Z"`
}
