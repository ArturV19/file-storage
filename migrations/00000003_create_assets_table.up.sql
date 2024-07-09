CREATE TABLE IF NOT EXISTS assets (
    name TEXT NOT NULL,
    uid BIGINT NOT NULL,
    data BYTEA NOT NULL, -- По-хорошему, должно быть в хранилище данных
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    original_name TEXT,
    content_type TEXT NOT NULL DEFAULT 'application/octet-stream',
    PRIMARY KEY (name, uid)
);

CREATE INDEX IF NOT EXISTS idx_assets_name ON assets (name);
CREATE INDEX IF NOT EXISTS idx_assets_uid ON assets (uid);