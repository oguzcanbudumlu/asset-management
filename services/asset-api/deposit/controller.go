package deposit

import (
	"asset-management/services/asset-api/dto"
	"github.com/gofiber/fiber/v2"
)

type Controller interface {
	Deposit(ctx *fiber.Ctx) error
}

type Request struct {
	WalletAddress string  `json:"wallet_address" example:"0x123abc456def"`
	Network       string  `json:"network" example:"Ethereum"`
	Amount        float64 `json:"amount" example:"100.50"`
}

type Response struct {
	NewBalance float64 `json:"new_balance" example:"1500.75"`
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
// @Param        depositRequest body Request true "Deposit request payload"
// @Success      200  {object}  Response
// @Failure      400  {object}  dto.ErrorResponse
// @Router       /deposit [post]
func (c *controller) Deposit(ctx *fiber.Ctx) error {
	var req Request
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Message: "Invalid request payload"})
	}

	newBalance, err := c.service.Deposit(req.WalletAddress, req.Network, req.Amount)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Message: err.Error()})
	}

	return ctx.JSON(Response{NewBalance: newBalance})
}
