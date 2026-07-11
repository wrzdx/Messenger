package chats_postgres_repository

import "messenger/internal/core/repository/postgres"

type ChatsRepository struct {
	db postgres.DB
}

func NewChatsRepository(db postgres.DB) *ChatsRepository {
	return &ChatsRepository{
		db: db,
	}
}
