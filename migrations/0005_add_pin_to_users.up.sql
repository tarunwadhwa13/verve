-- Migration: Add pin_hash to users table for 2FA
ALTER TABLE users
ADD COLUMN pin_hash VARCHAR(255);
