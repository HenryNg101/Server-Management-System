CREATE TABLE import_jobs (
    id UUID PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    file_path TEXT NOT NULL,
    status TEXT NOT NULL CHECK (
        status IN ('pending', 'processing', 'done', 'failed')
    ),
    processed_rows INT NOT NULL DEFAULT 0 CHECK (processed_rows >= 0),
    success_rows_count INT NOT NULL DEFAULT 0 CHECK (success_rows_count >= 0),
    failed_rows_count INT NOT NULL DEFAULT 0 CHECK (failed_rows_count >= 0),
    result_path TEXT,
    error TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);