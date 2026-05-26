CREATE TABLE servers (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    status BOOLEAN NOT NULL,
    ipv4_address INET NOT NULL,
    port INT NOT NULL,
    protocol TEXT DEFAULT 'tcp',
    created_at TIMESTAMP DEFAULT NOW(),
    last_updated TIMESTAMP DEFAULT NOW()
);

CREATE TABLE servers_users (
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    server_id INTEGER REFERENCES servers(id) ON DELETE CASCADE,
    user_role TEXT NOT NULL
)