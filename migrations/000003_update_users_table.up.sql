-- 1. Add columns with temporary defaults so existing rows don't break
ALTER TABLE users
ADD COLUMN password TEXT NOT NULL,
ADD COLUMN role TEXT NOT NULL;

-- 2. Add constraint for role
ALTER TABLE users
ADD CONSTRAINT users_role_check CHECK (role IN ('admin', 'user'));

-- 3. (Optional but recommended) Remove defaults so future inserts must be explicit
ALTER TABLE users
ALTER COLUMN password DROP DEFAULT,
ALTER COLUMN role DROP DEFAULT;