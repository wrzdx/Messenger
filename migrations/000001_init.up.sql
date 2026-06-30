CREATE TABLE users (
    id            BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    username      VARCHAR(32)  NOT NULL UNIQUE CHECK (char_length(username) BETWEEN 5 AND 32),
    first_name    VARCHAR(64)  NOT NULL CHECK (char_length(first_name) BETWEEN 1 AND 64),
    last_name     VARCHAR(64),
    created_at    TIMESTAMPTZ  NOT NULL,
    bio           VARCHAR(70),
    password_hash TEXT         NOT NULL
);