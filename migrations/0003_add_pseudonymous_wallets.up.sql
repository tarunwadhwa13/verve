-- Migration: Add pseudonymous wallet support and public key column

ALTER TABLE wallets
    ADD COLUMN is_pseudonymous BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN public_key TEXT;

-- Ensure wallet_id uniqueness for pseudonymous wallets
CREATE UNIQUE INDEX IF NOT EXISTS idx_wallets_pseudonymous_id ON wallets(id) WHERE is_pseudonymous;
