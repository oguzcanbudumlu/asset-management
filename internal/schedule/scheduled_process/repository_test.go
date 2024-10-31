package scheduled_process_test

import (
	"asset-management/internal/schedule/scheduled_process"
	"asset-management/services/asset-api/util"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPostgresProcessRepository_Process_EarlyExitCompletedTransaction(t *testing.T) {
	db, cleanup := util.SetupTestContainer(t)
	defer cleanup()

	repo := scheduled_process.NewProcessRepository(db)

	// Insert balance records for sender and receiver
	_, err := db.Exec(`INSERT INTO balance (wallet_address, network, balance) VALUES ($1, $2, $3)`,
		"wallet123", "mainnet", 200.0)
	assert.NoError(t, err)

	_, err = db.Exec(`INSERT INTO balance (wallet_address, network, balance) VALUES ($1, $2, $3)`,
		"wallet456", "mainnet", 100.0)
	assert.NoError(t, err)

	// Insert a completed scheduled transaction record
	_, err = db.Exec(`INSERT INTO scheduled_transactions (scheduled_transaction_id, from_wallet_address, to_wallet_address, network, amount, scheduled_time, status)
					  VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		123, "wallet123", "wallet456", "mainnet", 50.0, time.Now().Add(10*time.Minute), "COMPLETED")
	assert.NoError(t, err)

	// Process the transaction, expecting an early exit
	err = repo.Process(123)
	assert.NoError(t, err)

	// Verify the balances remain unchanged
	var senderBalance, receiverBalance float64
	err = db.QueryRow(`SELECT balance FROM balance WHERE wallet_address = $1 AND network = $2`, "wallet123", "mainnet").Scan(&senderBalance)
	assert.NoError(t, err)
	assert.Equal(t, 200.0, senderBalance)

	err = db.QueryRow(`SELECT balance FROM balance WHERE wallet_address = $1 AND network = $2`, "wallet456", "mainnet").Scan(&receiverBalance)
	assert.NoError(t, err)
	assert.Equal(t, 100.0, receiverBalance)

	// Verify the transaction status remains as COMPLETED
	var status string
	err = db.QueryRow(`SELECT status FROM scheduled_transactions WHERE scheduled_transaction_id = $1`, 123).Scan(&status)
	assert.NoError(t, err)
	assert.Equal(t, "COMPLETED", status)
}

func TestPostgresProcessRepository_Process_Success(t *testing.T) {
	db, cleanup := util.SetupTestContainer(t)
	defer cleanup()

	repo := scheduled_process.NewProcessRepository(db)

	// Insert balance records for sender and receiver
	_, err := db.Exec(`INSERT INTO balance (wallet_address, network, balance) VALUES ($1, $2, $3)`,
		"wallet123", "mainnet", 200.0)
	assert.NoError(t, err)

	_, err = db.Exec(`INSERT INTO balance (wallet_address, network, balance) VALUES ($1, $2, $3)`,
		"wallet456", "mainnet", 100.0)
	assert.NoError(t, err)

	// Insert scheduled transaction record with a scheduled_time
	_, err = db.Exec(`INSERT INTO scheduled_transactions (scheduled_transaction_id, from_wallet_address, to_wallet_address, network, amount, scheduled_time, status)
					  VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		123, "wallet123", "wallet456", "mainnet", 50.0, time.Now().Add(10*time.Minute), "PENDING")
	assert.NoError(t, err)

	// Process the transaction
	err = repo.Process(123)
	assert.NoError(t, err)

	// Verify the sender's and receiver's updated balances
	var senderBalance, receiverBalance float64
	err = db.QueryRow(`SELECT balance FROM balance WHERE wallet_address = $1 AND network = $2`, "wallet123", "mainnet").Scan(&senderBalance)
	assert.NoError(t, err)
	assert.Equal(t, 150.0, senderBalance)

	err = db.QueryRow(`SELECT balance FROM balance WHERE wallet_address = $1 AND network = $2`, "wallet456", "mainnet").Scan(&receiverBalance)
	assert.NoError(t, err)
	assert.Equal(t, 150.0, receiverBalance)

	// Verify transaction status update
	var status string
	err = db.QueryRow(`SELECT status FROM scheduled_transactions WHERE scheduled_transaction_id = $1`, 123).Scan(&status)
	assert.NoError(t, err)
	assert.Equal(t, "COMPLETED", status)
}

func TestPostgresProcessRepository_Process_InsufficientBalance(t *testing.T) {
	db, cleanup := util.SetupTestContainer(t)
	defer cleanup()

	repo := scheduled_process.NewProcessRepository(db)

	// Insert balance record for sender with insufficient balance
	_, err := db.Exec(`INSERT INTO balance (wallet_address, network, balance) VALUES ($1, $2, $3)`,
		"wallet123", "mainnet", 30.0)
	assert.NoError(t, err)

	// Insert scheduled transaction record with a scheduled_time
	_, err = db.Exec(`INSERT INTO scheduled_transactions (scheduled_transaction_id, from_wallet_address, to_wallet_address, network, amount, scheduled_time, status)
					  VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		123, "wallet123", "wallet456", "mainnet", 50.0, time.Now().Add(10*time.Minute), "PENDING")
	assert.NoError(t, err)

	// Process the transaction
	err = repo.Process(123)
	assert.Error(t, err)
	assert.Equal(t, "insufficient balance in sender's wallet", err.Error())

	// Verify that the status has not been changed
	var status string
	err = db.QueryRow(`SELECT status FROM scheduled_transactions WHERE scheduled_transaction_id = $1`, 123).Scan(&status)
	assert.NoError(t, err)
	assert.Equal(t, "PENDING", status)
}
