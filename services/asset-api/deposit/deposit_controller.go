package deposit

import (
	"github.com/gofiber/fiber/v2"
)

type DepositController interface {
	Deposit(ctx *fiber.Ctx) error
}

type depositController struct {
	service DepositService
}

func NewDepositController(service DepositService) DepositController {
	return &depositController{service: service}
}

// Deposit godoc
// @Summary      Deposit assets
// @Description  Deposits a specified amount into a wallet
// @Tags         deposit
// @Accept       json
// @Produce      json
// @Param        depositRequest body DepositRequest true "Deposit request payload"
// @Success      200  {object}  DepositResponse
// @Failure      400  {object}  ErrorResponse
// @Router       /deposit [post]
func (c *depositController) Deposit(ctx *fiber.Ctx) error {
	var req DepositRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: "Invalid request payload"})
	}

	newBalance, err := c.service.Deposit(req.WalletAddress, req.Network, req.Amount)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: err.Error()})
	}

	return ctx.JSON(DepositResponse{NewBalance: newBalance})
}
