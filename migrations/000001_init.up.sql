CREATE TYPE user_role AS ENUM ('member', 'admin', 'owner');
CREATE TYPE chat_type AS ENUM ('direct', 'group');

CREATE TABLE users (
    id            UUID         PRIMARY KEY,
    username      VARCHAR(32)  NOT NULL UNIQUE CHECK (char_length(username) BETWEEN 5 AND 32),
    first_name    VARCHAR(64)  NOT NULL CHECK (char_length(first_name) BETWEEN 1 AND 64),
    last_name     VARCHAR(64),
    created_at    TIMESTAMPTZ  NOT NULL,
    deleted_at    TIMESTAMPTZ,
    bio           VARCHAR(70),
    password_hash TEXT         NOT NULL
);

CREATE TABLE chats (
    id                UUID         PRIMARY KEY,
    type              chat_type    NOT NULL,
    name              VARCHAR(128),
    last_message_id   UUID,
    last_activity_at  TIMESTAMPTZ  NOT NULL,
    created_at        TIMESTAMPTZ  NOT NULL,

    CHECK(
        (type='group' AND name IS NOT NULL)
     OR (type='direct' AND name IS NULL)
    )
);

CREATE TABLE messages (
    id         UUID        PRIMARY KEY,
    chat_id    UUID        NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    sender_id  UUID        NOT NULL REFERENCES users(id),
    content    TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ
);

CREATE TABLE chat_participants (
    chat_id              UUID        NOT NULL REFERENCES chats(id),
    user_id              UUID        NOT NULL REFERENCES users(id),
    role                 user_role   NOT NULL,
    last_read_message_id UUID        REFERENCES messages(id),
    joined_at            TIMESTAMPTZ NOT NULL,

    PRIMARY KEY(chat_id, user_id)
);

CREATE TABLE directs (
    chat_id  UUID PRIMARY KEY REFERENCES chats(id),
    user1_id UUID NOT NULL REFERENCES users(id), 
    user2_id UUID NOT NULL REFERENCES users(id),

    UNIQUE(user1_id, user2_id)
    CHECK(user1_id <> user2_id) 
    CHECK (user1_id < user2_id)
);

ALTER TABLE chats
ADD CONSTRAINT fk_last_message
FOREIGN KEY (last_message_id)
REFERENCES messages(id);

CREATE INDEX idx_messages_chat_created
ON messages(chat_id, created_at DESC);

CREATE INDEX idx_chat_participants_user
ON chat_participants(user_id);

CREATE INDEX idx_chats_activity
ON chats(last_activity_at DESC);