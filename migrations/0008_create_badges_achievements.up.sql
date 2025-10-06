-- Migration: Create badges and achievements tables

-- Create badges table
CREATE TABLE badges (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT NOT NULL,
    icon_url TEXT,
    points INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by INTEGER REFERENCES users(id),
    is_active BOOLEAN DEFAULT TRUE
);

-- Create achievement_rules table for defining conditions
CREATE TABLE achievement_rules (
    id SERIAL PRIMARY KEY,
    badge_id INTEGER REFERENCES badges(id),
    rule_type VARCHAR(50) NOT NULL, -- e.g., 'transaction_count', 'transfer_amount', etc.
    condition_value JSONB NOT NULL, -- Store rule-specific conditions
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by INTEGER REFERENCES users(id),
    is_active BOOLEAN DEFAULT TRUE
);

-- Create user_badges table to track badge assignments
CREATE TABLE user_badges (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    badge_id INTEGER REFERENCES badges(id),
    awarded_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    awarded_by INTEGER REFERENCES users(id), -- can be NULL for system-awarded badges
    UNIQUE(user_id, badge_id)
);

-- Add badge management permissions
INSERT INTO permissions (name) VALUES 
    ('create_badge'),
    ('update_badge'),
    ('view_badges'),
    ('assign_badge');

-- Grant badge permissions to admin role
INSERT INTO role_permissions (role_id, permission_id) 
SELECT 
    r.id,
    p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'admin' 
AND p.name IN ('create_badge', 'update_badge', 'view_badges', 'assign_badge');
