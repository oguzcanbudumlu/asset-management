package scheduled_next_test

import (
	"asset-management/internal/schedule/scheduled_next"
	"asset-management/services/asset-api/util"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPostgresNextRepository_GetNextMinuteTransactions_Success(t *testing.T) {
	db, cleanup := util.SetupTestContainer(t)
	defer cleanup()

	repo := scheduled_next.NewNextRepository(db)

	// Insert a transaction scheduled for the next minute
	scheduledTime := time.Now().Add(30 * time.Second) // Set the scheduled time from code
	_, err := db.Exec(`
		INSERT INTO scheduled_transactions (from_wallet_address, to_wallet_address, network, amount, scheduled_time, status)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		"wallet123", "wallet456", "mainnet", 100.50, scheduledTime, "PENDING",
	)
	assert.NoError(t, err)

	transactions, err := repo.GetNextMinuteTransactions()

	assert.NoError(t, err)
	assert.NotEmpty(t, transactions)
	assert.Equal(t, "wallet123", transactions[0].FromWallet)
	assert.Equal(t, "wallet456", transactions[0].ToWallet)
}

func TestPostgresNextRepository_GetNextMinuteTransactions_EmptyResult(t *testing.T) {
	db, cleanup := util.SetupTestContainer(t)
	defer cleanup()

	repo := scheduled_next.NewNextRepository(db)

	// Ensure no transactions exist in the next minute
	transactions, err := repo.GetNextMinuteTransactions()

	assert.NoError(t, err)
	assert.Empty(t, transactions)
}
