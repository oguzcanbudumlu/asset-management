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

// GetWallets retrieves all wallets
// @Summary Get all wallets
// @Description Retrieve all wallets from the database
// @Produce json
// @Success 200 {array} Wallet
// @Failure 500 {string} string "Failed to retrieve wallets"
// @Router /wallet [get]
func (c *WalletController) GetWallets(ctx *fiber.Ctx) error {
	wallets, err := c.Service.GetWallets()
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).SendString("Failed to retrieve wallets")
	}
	return ctx.JSON(wallets)
}

// DeleteWallet deletes a wallet by address
// @Summary Delete a wallet
// @Description Delete a wallet by its address
// @Param address path string true "Wallet address"
// @Success 204 {string} string "No content"
// @Failure 404 {string} string "Wallet not found"
// @Router /wallet/{address} [delete]
func (c *WalletController) DeleteWallet(ctx *fiber.Ctx) error {
	address := ctx.Params("address")
	if err := c.Service.DeleteWallet(address); err != nil {
		return ctx.Status(http.StatusNotFound).SendString("Wallet not found")
	}
	return ctx.SendStatus(http.StatusNoContent)
}
