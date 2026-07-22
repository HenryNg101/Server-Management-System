CREATE TABLE agents (
    id SERIAL PRIMARY KEY,
    server_id INTEGER REFERENCES servers(id) ON DELETE CASCADE,
    api_key TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),

    instance_id TEXT NOT NULL,
    hostname TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'active',

    last_seen_at TIMESTAMP,
    CONSTRAINT unique_server_instance UNIQUE (server_id)
);