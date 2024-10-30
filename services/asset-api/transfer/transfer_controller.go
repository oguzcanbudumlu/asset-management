package transfer

import (
	"asset-management/services/asset-api/common/dto"
	"github.com/gofiber/fiber/v2"
)

type TransferController interface {
	Transfer(ctx *fiber.Ctx) error
}

type transferController struct {
	service TransferService
}

func NewDepositController(service TransferService) TransferController {
	return &transferController{service: service}
}

// Transfer godoc
// @Summary      Transfer assets
// @Description  Transfers a specified amount from a wallet into another
// @Tags         transfer
// @Accept       json
// @Produce      json
// @Param        transferRequest body TransferRequest true "Transfer request payload"
// @Success      200  "Transfer operation successful"
// @Failure      400  {object}  dto.ErrorResponse
// @Router       /transfer [post]
func (c *transferController) Transfer(ctx *fiber.Ctx) error {
	var req TransferRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Message: "Invalid request payload"})
	}

	err := c.service.Transfer(req.From, req.To, req.Network, req.Amount)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Message: err.Error()})
	}

	return ctx.SendStatus(fiber.StatusOK)
}
