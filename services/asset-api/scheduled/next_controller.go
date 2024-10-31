package scheduled

import (
	"asset-management/internal/schedule/scheduled_next"
	"github.com/gofiber/fiber/v2"
	"time"
)

type NextController struct {
	service scheduled_next.NextService
}

func NewNextController(service scheduled_next.NextService) *NextController {
	return &NextController{service: service}
}

type ScheduledTransaction struct {
	ID            int       `json:"id" example:"1"`                                // Transaction ID
	FromWallet    string    `json:"from_wallet" example:"wallet_123"`              // Sender's wallet address
	ToWallet      string    `json:"to_wallet" example:"wallet_456"`                // Recipient's wallet address
	Network       string    `json:"network" example:"Ethereum"`                    // Blockchain network (e.g., Ethereum)
	Amount        float64   `json:"amount" example:"250.75"`                       // Amount to be transferred
	ScheduledTime time.Time `json:"scheduled_time" example:"2024-10-30T15:04:05Z"` // Scheduled time for transaction
	Status        string    `json:"status" example:"PENDING"`                      // Transaction status (e.g., pending, completed)
	CreatedAt     time.Time `json:"created_at" example:"2024-10-29T10:15:00Z"`     // Time when the transaction was created
}

// GetNextMinuteTransactions godoc
// @Summary Get transactions scheduled for the next minute
// @Description Retrieve all transactions scheduled for the upcoming minute
// @Tags ScheduledTransaction
// @Accept  json
// @Produce  json
// @Success 200 {array} ScheduledTransaction
// @Failure 500 {object} map[string]interface{}
// @Router /scheduled-transaction/next [get]
func (c *NextController) GetNextMinuteTransactions(ctx *fiber.Ctx) error {
	transactions, err := c.service.GetNextMinuteTransactions()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve transactions"})
	}
	return ctx.JSON(transactions)
}
