package schedule

import (
	"github.com/gofiber/fiber/v2"
)

type NextController struct {
	service NextService
}

func NewNextController(service NextService) *NextController {
	return &NextController{service: service}
}

// GetNextMinuteTransactions godoc
// @Summary Get transactions scheduled for the next minute
// @Description Retrieve all transactions scheduled for the upcoming minute
// @Tags ScheduleTransaction
// @Accept  json
// @Produce  json
// @Success 200 {array} ScheduleTransaction
// @Failure 500 {object} map[string]interface{}
// @Router /scheduled-transaction/next-minute [get]
func (c *NextController) GetNextMinuteTransactions(ctx *fiber.Ctx) error {
	transactions, err := c.service.GetNextMinuteTransactions()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve transactions"})
	}
	return ctx.JSON(transactions)
}
