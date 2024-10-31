package util

import (
	sql2 "asset-management/internal/sql"
	"context"
	"database/sql"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
)

func SetupTestContainer(t *testing.T) (*sql.DB, func()) {
	t.Helper()
	ctx := context.Background()

	// Setting up PostgreSQL container
	req := testcontainers.ContainerRequest{
		Image:        "postgres:13",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpass",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}
	pgContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Failed to start container: %v", err)
	}

	// Retrieve the host and port for PostgreSQL
	host, err := pgContainer.Host(ctx)
	if err != nil {
		t.Fatalf("Failed to get container host: %v", err)
	}
	port, err := pgContainer.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("Failed to get container port: %v", err)
	}

	// Database connection setup
	dsn := "postgres://testuser:testpass@" + host + ":" + port.Port() + "/testdb?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}

	// Setup schema
	_, err = db.Exec(sql2.CreateBalanceTableSQL)
	assert.NoError(t, err)

	_, err = db.Exec(sql2.CreateScheduledTransactionsTable)
	assert.NoError(t, err)

	// Cleanup function to terminate the container
	cleanup := func() {
		db.Close()
		pgContainer.Terminate(ctx)
	}

	return db, cleanup
}

func InsertBalance(db *sql.DB, walletAddress, network string, balance float64) error {
	_, err := db.Exec(`
		INSERT INTO balance (wallet_address, network, balance) 
		VALUES ($1, $2, $3)`, walletAddress, network, balance)
	return err
}
