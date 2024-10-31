package sql

const CreateBalanceTableSQL = `
CREATE TABLE IF NOT EXISTS balance (
	wallet_address VARCHAR(255) NOT NULL,
	network VARCHAR(100) NOT NULL,
	balance NUMERIC(30, 10) NOT NULL DEFAULT 0,
	UNIQUE (wallet_address, network)
);
`

const CreateScheduledTransactionsTable = `
CREATE TABLE IF NOT EXISTS scheduled_transactions (
    scheduled_transaction_id SERIAL PRIMARY KEY,
    from_wallet_address VARCHAR(255) NOT NULL,
    to_wallet_address VARCHAR(255) NOT NULL,
    network VARCHAR(100) NOT NULL,
    amount NUMERIC(30, 10) NOT NULL CHECK (amount > 0),
    scheduled_time TIMESTAMP NOT NULL,
    status VARCHAR(50) DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'COMPLETED', 'FAILED')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`
