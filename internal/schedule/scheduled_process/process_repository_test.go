package scheduled_process

//
//import (
//	sql2 "asset-management/internal/sql"
//	"context"
//	"database/sql"
//	"fmt"
//	_ "github.com/lib/pq"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/require"
//	"github.com/testcontainers/testcontainers-go"
//	"github.com/testcontainers/testcontainers-go/wait"
//	"testing"
//	"time"
//)
//
//var db *sql.DB
//
//func SetupTestContainer(t *testing.T) (*sql.DB, func()) {
//	ctx := context.Background()
//
//	req := testcontainers.ContainerRequest{
//		Image:        "postgres:13",
//		ExposedPorts: []string{"5432/tcp"},
//		Env: map[string]string{
//			"POSTGRES_USER":     "testuser",
//			"POSTGRES_PASSWORD": "testpass",
//			"POSTGRES_DB":       "testdb",
//		},
//		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(2 * time.Second),
//	}
//
//	pgContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
//		ContainerRequest: req,
//		Started:          true,
//	})
//	require.NoError(t, err)
//
//	host, err := pgContainer.Host(ctx)
//	require.NoError(t, err)
//	port, err := pgContainer.MappedPort(ctx, "5432")
//	require.NoError(t, err)
//
//	dsn := fmt.Sprintf("postgres://testuser:testpass@%s:%s/testdb?sslmode=disable", host, port.Port())
//	db, err := sql.Open("postgres", dsn)
//	require.NoError(t, err)
//
//	err = db.Ping()
//	require.NoError(t, err)
//
//	setupSchema(t, db)
//
//	tearDown := func() {
//		db.Close()
//		pgContainer.Terminate(ctx)
//	}
//
//	return db, tearDown
//}
//
//func setupSchema(t *testing.T, db *sql.DB) {
//	_, balanceErr := db.Exec(sql2.CreateBalanceTableSQL)
//	_, scheduledErr := db.Exec(sql2.CreateScheduledTransactionsTable)
//	require.NoError(t, balanceErr)
//	require.NoError(t, scheduledErr)
//}
//
//func TestProcess_TransactionInitializationFailure(t *testing.T) {
//	db, tearDown := SetupTestContainer(t)
//	defer tearDown()
//
//	closeErr := db.Close()
//	require.NoError(t, closeErr)
//
//	repo := NewProcessRepository(db)
//	err := repo.Process(1)
//	assert.Error(t, err)
//	assert.Contains(t, err.Error(), "failed to begin transaction")
//}
//
//func TestProcess_FailedToDeductInsufficientFunds(t *testing.T) {
//	db, tearDown := SetupTestContainer(t)
//	defer tearDown()
//
//	_, err := db.Exec(`
//        INSERT INTO balance (wallet_address, network, balance) VALUES ('wallet1', 'mainnet', '50');
//        INSERT INTO balance (wallet_address, network, balance) VALUES ('wallet2', 'mainnet', '500');
//        INSERT INTO scheduled_transactions (
//            from_wallet_address,
//            to_wallet_address,
//            network,
//            amount,
//            scheduled_time,
//            status
//        )
//        VALUES ('wallet1', 'wallet2', 'mainnet', '100', NOW(), 'PENDING');
//    `)
//	require.NoError(t, err, "failed to insert initial test data")
//
//	repo := NewProcessRepository(db)
//
//	err = repo.Process(1)
//	assert.Error(t, err)
//	assert.Contains(t, err.Error(), "insufficient balance in sender's wallet")
//}
//
//func TestProcess_InsufficientBalance(t *testing.T) {
//	db, tearDown := SetupTestContainer(t)
//	defer tearDown()
//
//	ctx := context.Background()
//
//	_, err := db.ExecContext(ctx, `
//        INSERT INTO balance (wallet_address, network, balance)
//        VALUES ('wallet1', 'mainnet', '50');
//    `)
//	require.NoError(t, err, "failed to insert initial balance")
//
//	_, err = db.ExecContext(ctx, `
//        INSERT INTO scheduled_transactions (
//            from_wallet_address,
//            to_wallet_address,
//            network,
//            amount,
//            scheduled_time,
//            status
//        )
//        VALUES ('wallet1', 'wallet2', 'mainnet', '100', NOW(), 'PENDING');
//    `)
//	require.NoError(t, err, "failed to insert scheduled transaction")
//
//	repo := NewProcessRepository(db)
//
//	err = repo.Process(1)
//
//	require.Error(t, err, "expected an error due to insufficient balance")
//	require.Contains(t, err.Error(), "insufficient balance in sender's wallet")
//}
//
//func TestProcess_AddToReceiverBalance(t *testing.T) {
//	db, tearDown := SetupTestContainer(t)
//	defer tearDown()
//
//	ctx := context.Background()
//
//	_, err := db.ExecContext(ctx, `
//        INSERT INTO balance (wallet_address, network, balance)
//        VALUES ('sender_wallet', 'mainnet', '170');
//    `)
//	require.NoError(t, err, "failed to insert initial balance for sender")
//
//	_, err = db.ExecContext(ctx, `
//        INSERT INTO scheduled_transactions (
//            from_wallet_address,
//            to_wallet_address,
//            network,
//            amount,
//            scheduled_time,
//            status
//        )
//        VALUES ('sender_wallet', 'receiver_wallet', 'mainnet', '100', NOW(), 'PENDING')
//    `)
//	require.NoError(t, err, "failed to insert scheduled transaction")
//
//	_, err = db.ExecContext(ctx, `
//	   INSERT INTO balance (wallet_address, network, balance)
//	   VALUES ('receiver_wallet', 'mainnet', '40');
//	`)
//	require.NoError(t, err, "failed to insert initial balance for receiver")
//
//	repo := NewProcessRepository(db)
//
//	err = repo.Process(1)
//
//	var lastBalance float64
//	lastBalanceErr := db.QueryRowContext(ctx, `SELECT balance FROM balance
//        WHERE wallet_address = $1 AND network = $2`, "receiver_wallet", "mainnet").Scan(&lastBalance)
//	require.NoError(t, lastBalanceErr)
//
//	expectedBalance := float64(140)
//	assert.Equal(t, expectedBalance, lastBalance, "receiver's balance should not have increased")
//}
//
//func TestProcess_TransactionRowLockFailure(t *testing.T) {
//	//db, tearDown := SetupTestContainer(t)
//	//defer tearDown()
//	//
//	//_, err := db.Exec(`
//	//    INSERT INTO balance (wallet_address, network, balance) VALUES ('wallet1', 'mainnet', '500');
//	//    INSERT INTO balance (wallet_address, network, balance) VALUES ('wallet2', 'mainnet', '500');
//	//    INSERT INTO scheduled_transactions (
//	//        from_wallet_address,
//	//        to_wallet_address,
//	//        network,
//	//        amount,
//	//        scheduled_time,
//	//        status
//	//    )
//	//    VALUES ('wallet1', 'wallet2', 'mainnet', '100', NOW(), 'PENDING');
//	//`)
//	//require.NoError(t, err, "failed to insert initial test data")
//	//
//	//startSecondProcess := make(chan struct{})
//	//finishFirstProcess := make(chan struct{})
//	//
//	//go func() {
//	//	tx, err := db.Begin()
//	//	require.NoError(t, err, "failed to begin first transaction")
//	//
//	//	_, err = tx.Exec(`SELECT from_wallet_address, to_wallet_address, network, amount
//	//    FROM scheduled_transactions
//	//    WHERE scheduled_transaction_id = $1 FOR UPDATE`, 1)
//	//	require.NoError(t, err, "failed to acquire lock on transaction row")
//	//
//	//	close(startSecondProcess)
//	//
//	//	<-finishFirstProcess
//	//
//	//	err = tx.Commit()
//	//	require.NoError(t, err, "failed to commit first transaction")
//	//}()
//	//
//	//<-startSecondProcess
//	//
//	//repo := NewProcessRepository(db)
//	//err = repo.Process(1)
//	//assert.Error(t, err)
//	//assert.Contains(t, err.Error(), "failed to fetch scheduled transaction")
//	//
//	//close(finishFirstProcess)
//}
//
//func TestProcess_SenderRowLockFailure(t *testing.T) {
//	//db, tearDown := SetupTestContainer(t)
//	//defer tearDown()
//	//
//	//_, err := db.Exec(`
//	//    INSERT INTO balance (wallet_address, network, balance) VALUES ('wallet1', 'mainnet', '500');
//	//    INSERT INTO balance (wallet_address, network, balance) VALUES ('wallet2', 'mainnet', '500');
//	//    INSERT INTO scheduled_transactions (
//	//        from_wallet_address,
//	//        to_wallet_address,
//	//        network,
//	//        amount,
//	//        scheduled_time,
//	//        status
//	//    )
//	//    VALUES ('wallet1', 'wallet2', 'mainnet', '100', NOW(), 'PENDING');
//	//`)
//	//require.NoError(t, err, "failed to insert initial test data")
//	//
//	//startSecondProcess := make(chan struct{})
//	//finishFirstProcess := make(chan struct{})
//	//
//	//go func() {
//	//	tx, err := db.Begin()
//	//	require.NoError(t, err, "failed to begin first transaction")
//	//
//	//	_, err = tx.Exec(`SELECT balance FROM balance
//	//    WHERE wallet_address = $1 AND network = $2 FOR UPDATE`, "wallet1", "mainnet")
//	//	require.NoError(t, err, "failed to acquire lock on sender's balance")
//	//
//	//	close(startSecondProcess)
//	//
//	//	<-finishFirstProcess
//	//
//	//	err = tx.Commit()
//	//	require.NoError(t, err, "failed to commit first transaction")
//	//}()
//	//
//	//<-startSecondProcess
//	//
//	//repo := NewProcessRepository(db)
//	//err = repo.Process(1)
//	//assert.Error(t, err)
//	//assert.Contains(t, err.Error(), "failed to lock or retrieve sender's balance")
//	//
//	//close(finishFirstProcess)
//}
//
//func TestProcess_ReceiverRowLockFailure(t *testing.T) {
//	//db, tearDown := SetupTestContainer(t)
//	//defer tearDown()
//	//
//	//_, err := db.Exec(`
//	//    INSERT INTO balance (wallet_address, network, balance) VALUES ('wallet1', 'mainnet', '500');
//	//    INSERT INTO balance (wallet_address, network, balance) VALUES ('wallet2', 'mainnet', '500');
//	//    INSERT INTO scheduled_transactions (
//	//        from_wallet_address,
//	//        to_wallet_address,
//	//        network,
//	//        amount,
//	//        scheduled_time,
//	//        status
//	//    )
//	//    VALUES ('wallet1', 'wallet2', 'mainnet', '100', NOW(), 'PENDING');
//	//`)
//	//require.NoError(t, err, "failed to insert initial test data")
//	//
//	//startSecondProcess := make(chan struct{})
//	//finishFirstProcess := make(chan struct{})
//	//
//	//go func() {
//	//	tx, err := db.Begin()
//	//	require.NoError(t, err, "failed to begin first transaction")
//	//
//	//	_, err = tx.Exec(`SELECT balance FROM balance
//	//    WHERE wallet_address = $1 AND network = $2 FOR UPDATE`, "wallet2", "mainnet")
//	//	require.NoError(t, err, "failed to acquire lock on sender's balance")
//	//
//	//	close(startSecondProcess)
//	//
//	//	<-finishFirstProcess
//	//
//	//	err = tx.Commit()
//	//	require.NoError(t, err, "failed to commit first transaction")
//	//}()
//	//
//	//<-startSecondProcess
//	//
//	//repo := NewProcessRepository(db)
//	//err = repo.Process(1)
//	//assert.Error(t, err)
//	//assert.Contains(t, err.Error(), "failed to lock or retrieve receiver's balance")
//	//
//	//close(finishFirstProcess)
//}
