-- Migration: Add transactions and ledger_entries tables, and update wallets.balance to BIGINT

ALTER TABLE wallets
    ALTER COLUMN balance TYPE BIGINT;

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    sender_wallet_id INTEGER REFERENCES wallets(id),
    receiver_wallet_id INTEGER REFERENCES wallets(id),
    amount BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE ledger_entries (
    id SERIAL PRIMARY KEY,
    transaction_id INTEGER REFERENCES transactions(id),
    wallet_id INTEGER REFERENCES wallets(id),
    entry_type VARCHAR(10) NOT NULL CHECK (entry_type IN ('debit', 'credit')),
    amount BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    , extra TEXT
);
