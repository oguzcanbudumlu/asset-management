package main

import "github.com/gofiber/fiber/v2"

type DepositController interface {
	Deposit(ctx *fiber.Ctx) error
}

type depositController struct {
	service DepositService
}

func NewDepositController(service DepositService) DepositController {
	return &depositController{service: service}
}

func (c *depositController) Deposit(ctx *fiber.Ctx) error {
	var req DepositRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: "Invalid request payload"})
	}

	transactionID, newBalance, err := c.service.Deposit(req.WalletAddress, req.Network, req.Amount)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: err.Error()})
	}

	return ctx.JSON(DepositResponse{TransactionID: transactionID, NewBalance: newBalance})
}
