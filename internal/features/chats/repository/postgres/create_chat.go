package chats_postgres_repository

// func (r *ChatsRepository) createChat(
// 	ctx context.Context,
// 	chat domain.Chat,
// ) error {
// 	ctx, cancel := context.WithTimeout(ctx, r.timeout)
// 	defer cancel()
// 	db := postgres.GetExecutor(ctx, r.db)

// 	query := `
// INSERT INTO chats (id, type, last_activity_at, created_at)
// VALUES ($1, $2,$3,$4)
// 	RETURNING
// 		id,
// 		type,
// 		title,
// 		last_message_id,
// 		last_activity_at,
// 		created_at,
// 	`

// 	var model ChatModel
// 	err := db.QueryRow(ctx, query, chat.ID, chat.Type, chat.LastActivityAt, chat.CreatedAt).
// 		Scan(&model.ID, &model.Type, &model.Title, &model.LastMessageID, &model.LastActivityAt, &model.CreatedAt)

// 	if err != nil {
// 		return fmt.Errorf("create chat: %w", err)
// 	}

// 	return nil
// }
