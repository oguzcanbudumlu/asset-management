package withdraw

import (
	"asset-management/services/asset-api/common/dto"
	"github.com/gofiber/fiber/v2"
)

// WithdrawRequest represents the request payload for a withdrawal
type WithdrawRequest struct {
	WalletAddress string  `json:"wallet_address" example:"0x123abc456def"`
	Network       string  `json:"network" example:"Ethereum"`
	Amount        float64 `json:"amount" example:"100.50"`
}

// WithdrawResponse represents the response payload after a successful withdrawal
type WithdrawResponse struct {
	NewBalance float64 `json:"new_balance" example:"1500.75"`
}
type WithdrawController interface {
	Withdraw(ctx *fiber.Ctx) error
}

type withdrawController struct {
	service WithdrawService
}

func NewWithdrawController(service WithdrawService) WithdrawController {
	return &withdrawController{service: service}
}

// Withdraw godoc
// @Summary      Withdraw assets
// @Description  Withdraws a specified amount into a wallet
// @Tags         withdraw
// @Accept       json
// @Produce      json
// @Param        depositRequest body WithdrawRequest true "Withdraw request payload"
// @Success      200  "Withdraw operation successful"
// @Failure      400  {object}  dto.ErrorResponse
// @Router       /withdraw [post]
func (c *withdrawController) Withdraw(ctx *fiber.Ctx) error {
	var req WithdrawRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Message: "Invalid request payload"})
	}

	err := c.service.Withdraw(req.WalletAddress, req.Network, req.Amount)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Message: err.Error()})
	}

	return ctx.SendStatus(fiber.StatusOK)
}
