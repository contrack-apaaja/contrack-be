-- Manual database fix for role column
-- Run this SQL directly in your PostgreSQL database

-- Step 1: Add role column if it doesn't exist
DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'users' AND column_name = 'role'
    ) THEN
        ALTER TABLE users ADD COLUMN role VARCHAR(20) DEFAULT 'REGULAR';
    END IF;
END $$;

-- Step 2: Update existing users to have REGULAR role
UPDATE users SET role = 'REGULAR' WHERE role IS NULL OR role = '';

-- Step 3: Add check constraint
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_role_check;
ALTER TABLE users ADD CONSTRAINT users_role_check 
CHECK (role IN ('REGULAR', 'LEGAL', 'MANAGEMENT'));

-- Step 4: Create index
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);

-- Verify the changes
SELECT column_name, data_type, column_default 
FROM information_schema.columns 
WHERE table_name = 'users' AND column_name = 'role';
