-- Drop constraint first (Postgres requires it)
ALTER TABLE users
DROP CONSTRAINT IF EXISTS users_role_check;

-- Then drop columns
ALTER TABLE users
DROP COLUMN IF EXISTS password,
DROP COLUMN IF EXISTS role;