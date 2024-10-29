package main

import "github.com/gofiber/fiber/v2"

type WithdrawController interface {
	Deposit(ctx *fiber.Ctx) error
}

type withdrawController struct {
	service WithdrawService
}

func NewWithdrawController(service WithdrawService) WithdrawController {
	return &withdrawController{service: service}
}

func (c *withdrawController) Deposit(ctx *fiber.Ctx) error {
	return nil
}
