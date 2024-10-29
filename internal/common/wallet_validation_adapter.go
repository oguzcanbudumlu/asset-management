package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type WalletValidationAdapter interface {
	ValidateWallet(walletAddress, network string) (bool, error)
}

type walletValidationAdapter struct {
	baseURL string
}

func NewWalletValidationAdapter(baseURL string) WalletValidationAdapter {
	return &walletValidationAdapter{baseURL: baseURL}
}

func (a *walletValidationAdapter) ValidateWallet(walletAddress, network string) (bool, error) {
	// Create the request URL
	url := fmt.Sprintf("%s/wallets/validate?walletAddress=%s&network=%s", a.baseURL, walletAddress, network)

	// Make the HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// Check if the status code is 200 OK
	if resp.StatusCode != http.StatusOK {
		return false, errors.New("failed to validate wallet")
	}

	// Parse the JSON response (assuming it returns { "isValid": true/false })
	var result struct {
		IsValid bool `json:"isValid"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, err
	}

	return result.IsValid, nil
}
