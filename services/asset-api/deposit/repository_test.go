package deposit_test

import (
	"asset-management/services/asset-api/deposit"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRepository_Deposit(t *testing.T) {
	db, cleanup := SetupTestContainer(t)
	defer cleanup()

	repo := deposit.NewRepository(db)

	tests := []struct {
		name          string
		walletAddress string
		network       string
		amount        float64
		expectedErr   error
		expectedBal   float64
	}{
		{
			name:          "Initial deposit",
			walletAddress: "wallet_1",
			network:       "network_1",
			amount:        100.0,
			expectedErr:   nil,
			expectedBal:   100.0,
		},
		{
			name:          "Deposit to existing balance",
			walletAddress: "wallet_1",
			network:       "network_1",
			amount:        50.0,
			expectedErr:   nil,
			expectedBal:   150.0, // 100 + 50
		},
		{
			name:          "New deposit different wallet",
			walletAddress: "wallet_2",
			network:       "network_1",
			amount:        200.0,
			expectedErr:   nil,
			expectedBal:   200.0,
		},
		{
			name:          "Invalid deposit amount (negative)",
			walletAddress: "wallet_1",
			network:       "network_1",
			amount:        -100.0,
			expectedErr:   fmt.Errorf("deposit amount must be positive"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			balance, err := repo.Deposit(tt.walletAddress, tt.network, tt.amount)
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBal, balance)
			}
		})
	}
}
