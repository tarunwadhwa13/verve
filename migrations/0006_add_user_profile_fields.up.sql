-- Migration: Add profile fields and PIN configuration to users table
ALTER TABLE users
ADD COLUMN display_name VARCHAR(100),
ADD COLUMN profile_photo_url TEXT,
ADD COLUMN pin_required_for_transfer BOOLEAN NOT NULL DEFAULT TRUE;
