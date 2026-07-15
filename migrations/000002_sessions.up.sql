CREATE TABLE sessions (
    id               UUID,
    user_id          UUID        NOT NULL,
    current_token_id UUID        NOT NULL,
    last_used_at     TIMESTAMPTZ NOT NULL,
    created_at       TIMESTAMPTZ NOT NULL,
    expires_at       TIMESTAMPTZ NOT NULL,

    CONSTRAINT sessions_expires_after_created_check
    CHECK (expires_at > created_at),
    CONSTRAINT sessions_last_used_at_not_before_created_check
    CHECK (last_used_at >= created_at),
    CONSTRAINT sessions_last_used_at_not_after_expires_at
    CHECK (last_used_at <= expires_at),

    CONSTRAINT sessions_pkey PRIMARY KEY (id),
    CONSTRAINT sessions_user_id_fkey
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT sessions_current_token_id_key UNIQUE (current_token_id)
);

CREATE INDEX idx_sessions_user_id ON sessions(user_id);
