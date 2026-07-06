CREATE TABLE import_jobs (
    id UUID PRIMARY KEY,
    file_path TEXT NOT NULL,
    status TEXT NOT NULL CHECK (
        status IN ('pending', 'processing', 'done', 'failed')
    ),
    progress NUMERIC(5,2) NOT NULL DEFAULT 0 CHECK (
        progress >= 0.0 AND progress <= 100.0
    ),
    error TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);