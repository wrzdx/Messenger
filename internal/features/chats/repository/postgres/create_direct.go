package chats_postgres_repository

// func (r *ChatsRepository) createDirect(
// 	ctx context.Context,
// 	chat domain.Chat,
// 	user1 uuid.UUID,
// 	user2 uuid.UUID,
// ) error {
// ctx, cancel := context.WithTimeout(ctx, r.timeout)
// defer cancel()
// db := postgres.GetExecutor(ctx, r.db)

// query := `
// INSERT INTO directs (chat_id, user1_id, user2_id)
// VALUES ($1, $2, $3);
// `

// if _, err := db.Exec(ctx, query, chat.ID, user1, user2); err != nil {
// 	if errors.Is(err, postgres.ErrViolatesUnique) {
// 		return domain.AlreadyExistsErr(domain.ChatEntity, nil)
// 	}
// 	return fmt.Errorf("exec direct creation query: %w", err)
// }

// 	return nil
// }
