CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY DEFAULT encode(gen_random_bytes(16),'hex'),
    uid BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ip_address INET,
    expires_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_sessions_uid ON sessions (uid);
CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions (expires_at);