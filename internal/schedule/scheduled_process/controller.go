package scheduled_process

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type ProcessController struct {
	service ProcessService
}

func NewProcessController(service ProcessService) *ProcessController {
	return &ProcessController{service: service}
}

// Process godoc
// @Summary Process a scheduled transaction
// @Description Processes a scheduled transaction by its ID
// @Tags Transactions
// @Param id path int true "Transaction ID"
// @Success 200 {object} map[string]string "message": "Transaction processed successfully"
// @Failure 400 {object} map[string]string "error": "Invalid transaction ID"
// @Failure 500 {object} map[string]string "error": "Failed to process transaction"
// @Router /scheduled-transaction/{id}/process [post]
func (c *ProcessController) Process(ctx *fiber.Ctx) error {
	transactionIDParam := ctx.Params("id")
	transactionID, err := strconv.Atoi(transactionIDParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid transaction ID",
		})
	}

	if err := c.service.Process(transactionID); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Errorf("failed to process transaction: %w", err).Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Transaction processed successfully",
	})
}
