package wallet

import (
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
)

type ValidationAdapter interface {
	One(walletAddress, network string) error
	Both(from, to, network string) error
}

type walletValidationAdapter struct {
	baseURL string
}

func NewValidationAdapter(baseURL string) ValidationAdapter {
	return &walletValidationAdapter{baseURL: baseURL}
}

func (a *walletValidationAdapter) One(walletAddress, network string) error {
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

func (a *walletValidationAdapter) Both(from, to, network string) error {
	var g errgroup.Group

	g.Go(func() error {
		if err := a.One(from, network); err != nil {
			return errors.New("source wallet validation failed: " + err.Error())
		}
		return nil
	})

	g.Go(func() error {
		if err := a.One(to, network); err != nil {
			return errors.New("destination wallet validation failed: " + err.Error())
		}
		return nil
	})

	return g.Wait()
}
