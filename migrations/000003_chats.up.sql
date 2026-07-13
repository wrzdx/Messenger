CREATE TYPE user_role AS ENUM ('member', 'admin', 'owner');
CREATE TYPE chat_type AS ENUM ('direct', 'group');


CREATE TABLE chats (
    id                UUID         PRIMARY KEY,
    type              chat_type    NOT NULL,
    last_message_id   UUID,
    last_activity_at  TIMESTAMPTZ  NOT NULL,
    created_at        TIMESTAMPTZ  NOT NULL
);

CREATE TABLE groups (
    chat_id UUID PRIMARY KEY REFERENCES chats(id) ON DELETE CASCADE,
    title VARCHAR(128) NOT NULL CHECK (char_length(btrim(title)) BETWEEN 1 AND 128)
);

CREATE TABLE messages (
    id         UUID        PRIMARY KEY,
    chat_id    UUID        NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    sender_id  UUID        NOT NULL REFERENCES users(id),
    content    TEXT        NOT NULL CHECK (char_length(content) <= 4096),
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ
);

CREATE TABLE chat_participants (
    chat_id              UUID        NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    user_id              UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role                 user_role   NOT NULL,
    last_read_message_id UUID        REFERENCES messages(id),
    joined_at            TIMESTAMPTZ NOT NULL,

    PRIMARY KEY(chat_id, user_id)
);

CREATE TABLE directs (
    chat_id  UUID PRIMARY KEY REFERENCES chats(id) ON DELETE CASCADE,
    user1_id UUID NOT NULL REFERENCES users(id), 
    user2_id UUID NOT NULL REFERENCES users(id),

    UNIQUE(user1_id, user2_id),
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

CREATE INDEX idx_chat_type ON chats(type);
