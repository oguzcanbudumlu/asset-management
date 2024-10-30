package withdraw

import (
	"asset-management/services/asset-api/common/dto"
	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"
)

type WithdrawController interface {
	Withdraw(ctx *fiber.Ctx) error
}

type withdrawController struct {
	service WithdrawService
}

func NewWithdrawController(service WithdrawService) WithdrawController {
	return &withdrawController{service: service}
}

// WithdrawRequest represents the request payload for a withdrawal
type WithdrawRequest struct {
	WalletAddress string          `json:"wallet_address" example:"0x123abc456def"`
	Network       string          `json:"network" example:"Ethereum"`
	Amount        decimal.Decimal `json:"amount" example:"100.50"`
}

// WithdrawResponse represents the response payload after a successful withdrawal
type WithdrawResponse struct {
	NewBalance decimal.Decimal `json:"new_balance" example:"1500.75"`
}

// Withdraw godoc
// @Summary      Withdraw assets
// @Description  Withdraws a specified amount from a wallet
// @Tags         withdraw
// @Accept       json
// @Produce      json
// @Param        withdrawRequest body WithdrawRequest true "Withdraw request payload"
// @Success      200  {object}  WithdrawResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Router       /withdraw [post]
func (c *withdrawController) Withdraw(ctx *fiber.Ctx) error {
	var req WithdrawRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Message: "Invalid request payload"})
	}

	newBalance, err := c.service.Withdraw(req.WalletAddress, req.Network, req.Amount)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Message: err.Error()})
	}

	return ctx.JSON(WithdrawResponse{NewBalance: newBalance})
}
