-- Migration: Create transfers table to track transfer requests

CREATE TYPE transfer_status AS ENUM ('pending', 'completed', 'failed');

CREATE TABLE transfers (
    id SERIAL PRIMARY KEY,
    sender_wallet_id INTEGER REFERENCES wallets(id),
    receiver_wallet_id INTEGER REFERENCES wallets(id),
    amount BIGINT NOT NULL,
    status transfer_status NOT NULL DEFAULT 'pending',
    is_anonymous BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Add a trigger to update `updated_at` timestamp on row update
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_transfers_updated_at
BEFORE UPDATE ON transfers
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
