package wallet

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

func setupTestDB(t *testing.T) (*gorm.DB, func()) {
	ctx := context.Background()

	// Create PostgreSQL container
	req := testcontainers.ContainerRequest{
		Image:        "postgres:13",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       "testdb",
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpassword",
		},
		WaitingFor: wait.ForListeningPort("5432"),
	}

	pgContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Failed to start container: %s", err)
	}

	// Get the host and port of the container
	host, err := pgContainer.Host(ctx)
	if err != nil {
		t.Fatalf("Failed to get container host: %s", err)
	}
	port, err := pgContainer.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("Failed to get mapped port: %s", err)
	}

	dsn := fmt.Sprintf("host=%s port=%s user=testuser password=testpassword dbname=testdb sslmode=disable", host, port.Port())
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %s", err)
	}

	// Run migrations (if needed)
	if err := db.AutoMigrate(&Wallet{}, &WalletDeleted{}); err != nil {
		t.Fatalf("Failed to migrate database: %s", err)
	}

	// Cleanup function to stop the container
	return db, func() {
		_ = pgContainer.Terminate(ctx)
	}
}

func TestWalletRepository(t *testing.T) {
	db, teardown := setupTestDB(t)
	defer teardown()

	repo := NewWalletRepository(db)

	// Test CreateWallet
	wallet := &Wallet{Network: "Ethereum", Address: "0x123"}
	err := repo.CreateWallet(wallet)
	assert.NoError(t, err)

	// Test GetWallets
	wallets, err := repo.GetWallets()
	assert.NoError(t, err)
	assert.Len(t, wallets, 1)

	// Test GetWallet
	fetchedWallet, err := repo.GetWallet("Ethereum", "0x123")
	assert.NoError(t, err)
	assert.NotNil(t, fetchedWallet)
	assert.Equal(t, wallet.Address, fetchedWallet.Address)

	// Test DeleteWallet
	err = repo.DeleteWallet("Ethereum", "0x123")
	assert.NoError(t, err)

	// Check if the wallet was deleted
	fetchedWallet, err = repo.GetWallet("Ethereum", "0x123")
	assert.Error(t, err)
	assert.Nil(t, fetchedWallet)
}
