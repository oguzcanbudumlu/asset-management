package withdraw

import (
	"asset-management/services/asset-api/util"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRepository_Withdraw_Success(t *testing.T) {
	db, cleanup := util.SetupTestContainer(t)
	defer cleanup()

	repo := NewRepository(db)

	// Set up initial balance
	err := util.InsertBalance(db, "0x123abc456def", "Ethereum", 200.00)
	assert.NoError(t, err)

	// Act
	err = repo.Withdraw("0x123abc456def", "Ethereum", 100.50)

	// Assert
	assert.NoError(t, err)

	// Verify balance update
	var newBalance float64
	err = db.QueryRow(`
		SELECT balance FROM balance 
		WHERE wallet_address = $1 AND network = $2`, "0x123abc456def", "Ethereum").Scan(&newBalance)
	assert.NoError(t, err)
	assert.Equal(t, 99.50, newBalance)
}

func TestRepository_Withdraw_InsufficientBalance(t *testing.T) {
	db, cleanup := util.SetupTestContainer(t)
	defer cleanup()

	repo := NewRepository(db)

	// Set up initial balance
	err := util.InsertBalance(db, "0x123abc456def", "Ethereum", 50.00)
	assert.NoError(t, err)

	// Act
	err = repo.Withdraw("0x123abc456def", "Ethereum", 100.50)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "insufficient balance", err.Error())
}

func TestRepository_Withdraw_WalletNotFound(t *testing.T) {
	db, cleanup := util.SetupTestContainer(t)
	defer cleanup()

	repo := NewRepository(db)

	// Act
	err := repo.Withdraw("0x123abc456def", "Ethereum", 100.50)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "wallet not found", err.Error())
}
