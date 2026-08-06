-- Development-only bootstrap users.
-- Password for both accounts: password
-- The value is a bcrypt hash accepted by the auth service.
-- The WHERE clauses make the script safe to run more than once.

INSERT INTO users (name, email, password, role)
SELECT 'Development Administrator', 'admin@example.com',
       '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
       'admin'
WHERE NOT EXISTS (
    SELECT 1 FROM users WHERE email = 'admin@example.com'
);

INSERT INTO users (name, email, password, role)
SELECT 'Development User', 'user@example.com',
       '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
       'user'
WHERE NOT EXISTS (
    SELECT 1 FROM users WHERE email = 'user@example.com'
);
