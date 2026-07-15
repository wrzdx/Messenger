CREATE TABLE users (
    id            UUID,
    username      VARCHAR(32)  NOT NULL,
    first_name    VARCHAR(64)  NOT NULL,
    last_name     VARCHAR(64),
    created_at    TIMESTAMPTZ  NOT NULL,
    deleted_at    TIMESTAMPTZ,
    bio           VARCHAR(70),
    password_hash TEXT         NOT NULL,

    CONSTRAINT users_username_check 
        CHECK (username ~ '^[a-zA-Z0-9_]{5,32}$'),
    CONSTRAINT users_first_name_check 
        CHECK (char_length(first_name) BETWEEN 1 AND 64),
    CONSTRAINT users_deleted_at_after_created_at_check
        CHECK (deleted_at IS NULL OR deleted_at >= created_at),
    CONSTRAINT users_password_hash_not_empty_check
        CHECK (password_hash <> ''),
    CONSTRAINT users_last_name_not_blank_check
        CHECK (last_name IS NULL OR btrim(last_name) <> ''),
    CONSTRAINT users_bio_not_blank_check
        CHECK (bio IS NULL OR btrim(bio) <> ''),

    CONSTRAINT users_pkey PRIMARY KEY (id)
);

CREATE UNIQUE INDEX users_username_lower_uidx
    ON users (lower(username));
