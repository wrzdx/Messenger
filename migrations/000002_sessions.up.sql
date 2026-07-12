CREATE TABLE sessions (
    id               UUID        PRIMARY KEY,
    user_id          UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    current_token_id UUID        NOT NULL UNIQUE,
    last_used_at     TIMESTAMPTZ NOT NULL,
    created_at       TIMESTAMPTZ NOT NULL,
    expires_at       TIMESTAMPTZ NOT NULL,

    CHECK (expires_at > created_at),
    CHECK (last_used_at >= created_at),
    CHECK (last_used_at <= expires_at)
);

CREATE INDEX idx_sessions_user_id ON sessions(user_id);
