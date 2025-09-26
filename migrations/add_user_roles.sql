-- Migration: Add role column to users table
-- File: migrations/add_user_roles.sql

ALTER TABLE users ADD COLUMN IF NOT EXISTS role VARCHAR(20) DEFAULT 'user';

-- Add constraint to ensure only valid roles
ALTER TABLE users ADD CONSTRAINT check_user_role 
CHECK (role IN ('user', 'legal', 'management'));

-- Create index for faster role-based queries
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);

-- Update existing users to have 'user' role if NULL
UPDATE users SET role = 'user' WHERE role IS NULL;

-- Insert some sample roles for testing (optional)
INSERT INTO users (email, password_hash, name, role, created_at) VALUES 
('legal@company.com', '$2a$10$hash', 'Legal User', 'legal', NOW()),
('manager@company.com', '$2a$10$hash', 'Management User', 'management', NOW())
ON CONFLICT (email) DO NOTHING;