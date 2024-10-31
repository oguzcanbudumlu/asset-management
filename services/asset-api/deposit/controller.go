package deposit

import (
	"asset-management/services/asset-api/common/dto"
	"github.com/gofiber/fiber/v2"
)

type Controller interface {
	Deposit(ctx *fiber.Ctx) error
}

type controller struct {
	service Service
}

func NewController(service Service) Controller {
	return &controller{service: service}
}

// Deposit godoc
// @Summary      Deposit assets
// @Description  Deposits a specified amount into a wallet
// @Tags         deposit
// @Accept       json
// @Produce      json
// @Param        depositRequest body DepositRequest true "Deposit request payload"
// @Success      200  {object}  DepositResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Router       /deposit [post]
func (c *controller) Deposit(ctx *fiber.Ctx) error {
	var req DepositRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Message: "Invalid request payload"})
	}

	newBalance, err := c.service.Deposit(req.WalletAddress, req.Network, req.Amount)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Message: err.Error()})
	}

	return ctx.JSON(DepositResponse{NewBalance: newBalance})
}
