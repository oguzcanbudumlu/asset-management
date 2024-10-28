package main

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
)

// WalletController struct
type WalletController struct {
	Service *WalletService
}

// NewWalletController creates a new WalletController
func NewWalletController(service *WalletService) *WalletController {
	return &WalletController{Service: service}
}

// CreateWallet creates a new wallet
// @Summary Create a wallet
// @Description Create a new wallet with an address and network
// @Accept json
// @Produce json
// @Param wallet body Wallet true "Wallet data"
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Failed to create wallet"
// @Router /wallet [post]
func (c *WalletController) CreateWallet(ctx *fiber.Ctx) error {
	var wallet Wallet
	if err := ctx.BodyParser(&wallet); err != nil {
		return ctx.Status(http.StatusBadRequest).SendString("Invalid request")
	}

	if err := c.Service.CreateWallet(&wallet); err != nil {
		return ctx.Status(http.StatusInternalServerError).SendString("Failed to create wallet")
	}

	return ctx.Status(http.StatusCreated).JSON(wallet)
}

// GetWallet retrieves a wallet by network and address
// @Summary Get a wallet by network and address
// @Description Retrieve a wallet by its network and address
// @Param network path string true "Wallet network"
// @Param address path string true "Wallet address"
// @Success 200 {object} Wallet
// @Failure 404 {string} string "Wallet not found"
// @Failure 500 {string} string "Failed to retrieve wallet"
// @Router /wallet/{network}/{address} [get]
func (c *WalletController) GetWallet(ctx *fiber.Ctx) error {
	network := ctx.Params("network")
	address := ctx.Params("address")

	wallet, err := c.Service.GetWallet(network, address)
	if err != nil {
		return ctx.Status(http.StatusNotFound).SendString("Wallet not found")
	}
	return ctx.JSON(wallet)
}

// DeleteWallet deletes a wallet by address and network
// @Summary Delete a wallet
// @Description Delete a wallet by its address and network
// @Param address path string true "Wallet address"
// @Param network path string true "Wallet network"
// @Success 204 {string} string "No content"
// @Failure 404 {string} string "Wallet not found"
// @Failure 500 {string} string "Failed to delete wallet"
// @Router /wallet/{network}/{address} [delete]
func (c *WalletController) DeleteWallet(ctx *fiber.Ctx) error {
	network := ctx.Params("network")
	address := ctx.Params("address")

	if err := c.Service.DeleteWallet(network, address); err != nil {
		return ctx.Status(http.StatusNotFound).SendString("Wallet not found")
	}
	return ctx.SendStatus(http.StatusNoContent)
}
