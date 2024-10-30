package transaction

import (
	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"
	"time"
)

type CreateController struct {
	service CreateService
}

func NewCreateController(service CreateService) *CreateController {
	return &CreateController{service: service}
}

type Request struct {
	From          string          `json:"from" example:"wallet123"`
	To            string          `json:"to" example:"wallet456"`
	Network       string          `json:"network" example:"mainnet"`
	Amount        decimal.Decimal `json:"amount" example:"100.50"`
	ScheduledTime string          `json:"scheduled_time" example:"2023-12-31T12:00:00Z"`
}

// Create godoc
// @Summary      Create a new scheduled transaction
// @Description  Schedules a new transaction to be executed at a specified future time
// @Tags         ScheduleTransaction
// @Accept       json
// @Produce      json
// @Param        transaction body Request true "Schedule Transfer request payload"
// @Success      201  {object}  map[string]int "Created transaction ID"  example: {"transaction_id": 123}
// @Failure      400  {object}  map[string]string "Invalid request payload or scheduled time format" example: {"error": "Invalid scheduled time format"}
// @Failure      500  {object}  map[string]string "Failed to create scheduled transaction" example: {"error": "Failed to create scheduled transaction"}
// @Router       /scheduled-transaction [post]
func (c *CreateController) Create(ctx *fiber.Ctx) error {
	var req Request
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	scheduledTime, err := time.Parse(time.RFC3339, req.ScheduledTime)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid scheduled time format"})
	}

	id, err := c.service.Create(req.From, req.To, req.Network, req.Amount, scheduledTime)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"transaction_id": id})
}
