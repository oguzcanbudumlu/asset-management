package withdraw

import (
	"asset-management/services/asset-api/common/dto"
	"github.com/gofiber/fiber/v2"
)

type Request struct {
	WalletAddress string  `json:"wallet_address" example:"0x123abc456def"`
	Network       string  `json:"network" example:"Ethereum"`
	Amount        float64 `json:"amount" example:"100.50"`
}

type Response struct {
	NewBalance float64 `json:"new_balance" example:"1500.75"`
}
type Controller interface {
	Withdraw(ctx *fiber.Ctx) error
}

type controller struct {
	service Service
}

func NewController(service Service) Controller {
	return &controller{service: service}
}

// Withdraw godoc
// @Summary      Withdraw assets
// @Description  Withdraws a specified amount into a wallet
// @Tags         withdraw
// @Accept       json
// @Produce      json
// @Param        depositRequest body Request true "Withdraw request payload"
// @Success      200  "Withdraw operation successful"
// @Failure      400  {object}  dto.ErrorResponse
// @Router       /withdraw [post]
func (c *controller) Withdraw(ctx *fiber.Ctx) error {
	var req Request
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Message: "Invalid request payload"})
	}

	err := c.service.Withdraw(req.WalletAddress, req.Network, req.Amount)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Message: err.Error()})
	}

	return ctx.SendStatus(fiber.StatusOK)
}
