package deposit_test

import (
	"asset-management/services/asset-api/deposit"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
	"time"
)

func TestRepository_Deposit(t *testing.T) {
	// Setup the test container and database schema
	db, cleanup, err := SetupTestContainer(t)
	if err != nil {
		t.Fatalf("failed to set up test container: %v", err)
	}
	defer cleanup()

	repo := deposit.NewRepository(db)

	tests := []struct {
		name          string
		walletAddress string
		network       string
		amount        float64
		expectedErr   error
		expectedBal   float64
		setup         func() // optional setup for each test case
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
			// Run setup function if defined
			if tt.setup != nil {
				tt.setup()
			}

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

// SetupTestContainer initializes the PostgreSQL test container and sets up the database schema
func SetupTestContainer(t *testing.T) (*sql.DB, func(), error) {
	ctx := context.Background()

	// Setup the test container and return the database connection and cleanup function
	db, cleanup, err := InitializeContainerAndDB(ctx, t) // assumes a function that initializes the container and db
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize container and db: %w", err)
	}

	// Setup schema (e.g., create balance table)
	if err := setupSchema(db); err != nil {
		cleanup()
		return nil, nil, err
	}

	return db, cleanup, nil
}

// setupSchema initializes the balance table in the database for testing
func setupSchema(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS balance (
			wallet_address VARCHAR(255) NOT NULL,
			network VARCHAR(255) NOT NULL,
			balance NUMERIC(30, 10) NOT NULL CHECK (balance >= 0),
			PRIMARY KEY (wallet_address, network)
		);
	`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}
	return nil
}

func InitializeContainerAndDB(ctx context.Context, t *testing.T) (*sql.DB, func(), error) {
	// Define the PostgreSQL container request
	req := testcontainers.ContainerRequest{
		Image:        "postgres:13",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpassword",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(60 * time.Second),
	}

	// Start the container
	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start postgres container: %w", err)
	}

	// Get the container's mapped port
	host, err := postgresContainer.Host(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get container host: %w", err)
	}
	port, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get mapped port: %w", err)
	}

	// Create the database connection
	dsn := fmt.Sprintf("postgres://testuser:testpassword@%s:%s/testdb?sslmode=disable", host, port.Port())
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		postgresContainer.Terminate(ctx)
		return nil, nil, fmt.Errorf("failed to open db connection: %w", err)
	}

	// Check the database connection
	if err := db.Ping(); err != nil {
		postgresContainer.Terminate(ctx)
		return nil, nil, fmt.Errorf("failed to ping db: %w", err)
	}

	// Cleanup function to terminate the container
	cleanup := func() {
		db.Close()
		postgresContainer.Terminate(ctx)
	}

	return db, cleanup, nil
}
