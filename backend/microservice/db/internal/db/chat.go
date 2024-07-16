package db

import (
	"InterestingChats/backend/microservice/db/internal/models"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type ChatService struct {
	Db *sql.DB
}

func NewChatService(db *sql.DB) *ChatService {
	return &ChatService{
		Db: db,
	}
}

func (cs *ChatService) GetChatInfo(chatName string) (*models.Chat, error) {
	chat := &models.Chat{
		Members: make(map[int]models.User),
	}

	err := cs.Db.QueryRowContext(context.Background(), "SELECT id, name FROM chats WHERE name = $1", chatName).Scan(&chat.ID, &chat.ChatName)
	if err != nil {
		return nil, fmt.Errorf("error select chat details from database: %w", err)
	}
	log.Printf("Chat details: ID=%d, Name=%s\n", chat.ID, chat.ChatName)

	rows, err := cs.Db.QueryContext(context.Background(), "SELECT u.id, u.username, u.email FROM chat_members cm JOIN users u ON cm.user_id = u.id WHERE cm.chat_id = $1", chat.ID)
	if err != nil {
		return nil, fmt.Errorf("error select chat detiels from database: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email); err != nil {
			return nil, fmt.Errorf("error scan query result(users) from database: %w", err)
		}
		chat.Members[user.ID] = user
		log.Printf("Added member: ID=%d, Username=%s\n, email=%s\n", user.ID, user.Username, user.Email)
	}
	rows, err = cs.Db.QueryContext(context.Background(), "SELECT id, body, created_at, user_id FROM messages WHERE chat_id = $1 ORDER BY created_at", chat.ID)
	if err != nil {
		return nil, fmt.Errorf("error select messages detiels from database: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(&msg.ID, &msg.Body, &msg.Time, &msg.UserID); err != nil {
			return nil, fmt.Errorf("error scan query result(messages) from database: %w", err)
		}
		chat.Messages = append(chat.Messages, msg)
		log.Printf("Added message: ID=%d, Body=%s, Time=%s, UserID=%d\n", msg.ID, msg.Body, msg.Time, msg.UserID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over messages rows: %w", err)
	}

	log.Printf("Final chat details: ID=%d, Name=%s, Members=%v, Messages=%v\n", chat.ID, chat.ChatName, chat.Members, chat.Messages)
	return chat, nil

}

func (cs *ChatService) CreateChat(chat *models.Chat) error {
	err := cs.Db.QueryRowContext(context.Background(), "INSERT INTO chats(name, created_at) VALUES($1, $2) RETURNING id", chat.ChatName, time.Now()).Scan(&chat.ID)
	if err != nil {
		return fmt.Errorf("failed insert into database(chat table): %w", err)
	}

	insertMemberStmt, err := cs.Db.PrepareContext(context.Background(), "INSERT INTO chat_members(chat_id, user_id)")
	if err != nil {
		return fmt.Errorf("failed prepare insert into database(chat_members table): %w", err)
	}
	defer insertMemberStmt.Close()

	for _, user := range chat.Members {
		_, err := insertMemberStmt.ExecContext(context.Background(), chat.ID, user.ID)
		if err != nil {
			return fmt.Errorf("failed cycle insert into database(chat_members table): %w", err)
		}
	}
	log.Println(chat.Members, chat.ID, chat.Messages, chat.ChatName)
	return nil
}
