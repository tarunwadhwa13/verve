-- Migration: Create treasury user, role, and wallet
-- This migration seeds the database with the necessary components for the treasury system.

-- 1. Create a non-privileged 'treasury' role
INSERT INTO roles (name) VALUES ('treasury');

-- 2. Create the dedicated treasury user
-- In a real production environment, this password should be set via a secure mechanism.
INSERT INTO users (username, password_hash) VALUES ('treasury@system.local', '$2a$10$D8.v4.gY9.eN9.pG.E.L3uF6.X8.Y.pG.E.L3uF6.X8.Y.pG.E.L3'); -- password is "password"

-- 3. Assign the 'treasury' role to the new user
INSERT INTO user_roles (user_id, role_id) VALUES
((SELECT id FROM users WHERE username = 'treasury@system.local'), (SELECT id FROM roles WHERE name = 'treasury'));

-- 4. Create the treasury wallet and link it to the treasury user
INSERT INTO wallets (user_id, currency, balance) VALUES
((SELECT id FROM users WHERE username = 'treasury@system.local'), 'USD', 0);
