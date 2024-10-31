package transaction_test

import (
	"asset-management/internal/schedule"
	"asset-management/services/asset-api/transaction"
	"asset-management/services/asset-api/util"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPostgresCreateRepository_Create_Success(t *testing.T) {
	db, cleanup := util.SetupTestContainer(t)
	defer cleanup()

	repo := transaction.NewCreateRepository(db)

	tx := &schedule.ScheduleTransaction{
		FromWallet:    "wallet123",
		ToWallet:      "wallet456",
		Network:       "mainnet",
		Amount:        100.50,
		ScheduledTime: time.Now().Add(24 * time.Hour),
		Status:        schedule.StatusPending,
	}

	id, err := repo.Create(tx)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, id)
}

func TestPostgresCreateRepository_Create_Failure(t *testing.T) {
	db, cleanup := util.SetupTestContainer(t)
	defer cleanup()

	repo := transaction.NewCreateRepository(db)

	tx := &schedule.ScheduleTransaction{
		FromWallet:    "",
		ToWallet:      "wallet456",
		Network:       "mainnet",
		Amount:        0,
		ScheduledTime: time.Now().Add(24 * time.Hour),
		Status:        schedule.StatusPending,
	}

	id, err := repo.Create(tx)
	assert.Error(t, err)
	assert.Equal(t, 0, id)
}
