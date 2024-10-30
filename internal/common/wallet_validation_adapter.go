package common

import (
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
)

type WalletValidationAdapter interface {
	ValidateWallet(walletAddress, network string) error
	ValidateBoth(from, to, network string) error
}

type walletValidationAdapter struct {
	baseURL string
}

func NewWalletValidationAdapter(baseURL string) WalletValidationAdapter {
	return &walletValidationAdapter{baseURL: baseURL}
}

func (a *walletValidationAdapter) ValidateWallet(walletAddress, network string) error {
	// Create the request URL
	url := fmt.Sprintf("%s/wallet/%s/%s", a.baseURL, network, walletAddress)

	// Make the HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check if the status code is 200 OK
	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to validate wallet")
	}

	return nil
}

func (a *walletValidationAdapter) ValidateBoth(from, to, network string) error {
	var g errgroup.Group

	g.Go(func() error {
		if err := a.ValidateWallet(from, network); err != nil {
			return errors.New("source wallet validation failed: " + err.Error())
		}
		return nil
	})

	g.Go(func() error {
		if err := a.ValidateWallet(to, network); err != nil {
			return errors.New("destination wallet validation failed: " + err.Error())
		}
		return nil
	})

	return g.Wait()
}
