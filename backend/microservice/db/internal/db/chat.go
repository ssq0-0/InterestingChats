package db

import (
	"InterestingChats/backend/microservice/db/internal/models"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/lib/pq"
)

type ChatService struct {
	Db *sql.DB
}

func NewChatService(dbConnection *sql.DB) *ChatService {
	if dbConnection == nil {
		log.Fatal("dbConnection is nil")
	}
	return &ChatService{
		Db: dbConnection,
	}
}

func (cs *ChatService) GetChatInfo(chatName string) (*models.Chat, error) {
	chat := &models.Chat{
		Members: make(map[int]models.User),
	}

	tx, err := cs.Db.BeginTx(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create 'begin transaction': %w", err)
	}

	defer func() {
		if p := recover(); p != nil || err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("failed to rollback: %v", err)
			}
		}
	}()

	err = tx.QueryRowContext(context.Background(), "SELECT id, name FROM chats WHERE name = $1", chatName).Scan(&chat.ID, &chat.ChatName)
	if err != nil {
		return nil, fmt.Errorf("error select chat details from database: %w", err)
	}
	log.Printf("Chat details: ID=%d, Name=%s", chat.ID, chat.ChatName)

	memberRows, err := tx.QueryContext(context.Background(), "SELECT u.id, u.username, u.email FROM chat_members cm JOIN users u ON cm.user_id = u.id WHERE cm.chat_id = $1", chat.ID)
	if err != nil {
		return nil, fmt.Errorf("error select chat members from database: %w", err)
	}
	defer memberRows.Close()

	for memberRows.Next() {
		var user models.User
		if err := memberRows.Scan(&user.ID, &user.Username, &user.Email); err != nil {
			return nil, fmt.Errorf("error scan query result(users) from database: %w", err)
		}
		chat.Members[user.ID] = user
		log.Printf("Added member: ID=%d, Username=%s\n, Email=%s\n", user.ID, user.Username, user.Email)
	}

	if memberRows.Err(); err != nil {
		return nil, fmt.Errorf("failed to ittarating over chat members rows: %w", err)
	}

	messageRows, err := tx.QueryContext(context.Background(), "SELECT id, body, created_at, user_id FROM messages WHERE chat_id = $1 ORDER BY created_at", chat.ID)
	if err != nil {
		return nil, fmt.Errorf("error select messages from database: %w", err)
	}
	defer messageRows.Close()

	for messageRows.Next() {
		var msg models.Message
		if err := messageRows.Scan(&msg.ID, &msg.Body, &msg.Time, &msg.UserID); err != nil {
			return nil, fmt.Errorf("error scan query result(messages) from database: %w", err)
		}
		chat.Messages = append(chat.Messages, msg)
		log.Printf("Added message: ID=%d, Body=%s, Time=%s, UserID=%d", msg.ID, msg.Body, msg.Time, msg.UserID)
	}

	if err := messageRows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over messages rows: %w", err)
	}

	if err := tx.Commit(); err != nil {
		log.Printf("failed to commit transaction for get chat info: %v", err)
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return chat, nil

}

func (cs *ChatService) CreateChat(chat *models.Chat) (int, error) {
	tx, err := cs.Db.BeginTx(context.Background(), nil)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil || err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("failed to rollback: %v", err)
			}
		}
	}()

	err = tx.QueryRowContext(context.Background(), "INSERT INTO chats(name, created_at) VALUES($1, $2) RETURNING id", chat.ChatName, time.Now()).Scan(&chat.ID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return 0, fmt.Errorf("failed to create chat with: %s name(already exist.)", chat.ChatName)
		}
		return 0, fmt.Errorf("failed insert into database(chat table): %w \n ", err)
	}
	log.Printf("chat id: %d", chat.ID)

	insertMemberStmt, err := tx.PrepareContext(context.Background(), "INSERT INTO chat_members(chat_id, user_id) VALUES($1, $2)")
	if err != nil {
		return 0, fmt.Errorf("failed prepare insert into database(chat_members table): %w", err)
	}
	defer insertMemberStmt.Close()

	for _, user := range chat.Members {
		_, err := insertMemberStmt.ExecContext(context.Background(), chat.ID, user.ID)
		if err != nil {
			return 0, fmt.Errorf("failed to insert into database(chat_members table): %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("failed to commit transaction: %v", err)
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return chat.ID, err
}

func (cs *ChatService) DeleteChat(chatID int) error {
	tx, err := cs.Db.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("failed create transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil || err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("failed to rollback transaction: %v", rbErr)
			}
		}
	}()

	_, err = tx.ExecContext(context.Background(), "DELETE FROM chats WHERE id = $1", chatID)
	if err != nil {
		log.Printf("failed to exec: %v", err)
		return nil
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed commit transaction: %w", err)
	}

	return nil
}

func (cs *ChatService) CheckAuthor(email string, chatID int) (bool, error) {
	var result string
	err := cs.Db.QueryRowContext(context.Background(), "SELECT u.email FROM chats c JOIN users u ON c.author_id = u.id WHERE c.id=$1", chatID).Scan(&result)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("no author found for chat id: %d", chatID)
			return false, nil
		}
		log.Printf("failed exec from db: %v", err)
		return false, err
	}

	if result != email {
		log.Printf("email mismatch: expected %s, got %s", email, result)
		return false, nil
	}

	return true, nil
}

func (cs *ChatService) DeleteMember(chatID, userID int) error {
	tx, err := cs.Db.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("failed create tx: %w", err)
	}
	defer func() {
		if p := recover(); p != nil || err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("failed to rollback: %v", err)
			}
		}
	}()

	var exists bool
	err = tx.QueryRowContext(context.Background(), "SELECT EXISTS(SELECT 1 FROM chat_members WHERE chat_id=$1 AND user_id=$2)", chatID, userID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if user exists in chat: %v", err)
	}
	if !exists {
		log.Printf("user %d is not member of chat %d", userID, chatID)
		return fmt.Errorf("user %d is not member of chat %d", userID, chatID)
	}

	_, err = tx.ExecContext(context.Background(), "DELETE FROM chat_members WHERE chat_id = $1 AND user_id=$2", chatID, userID)
	if err != nil {
		return fmt.Errorf("failed to execute delete query: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to comit delete tx: %w", err)
	}

	return nil
}

func (cs *ChatService) AddMember(chatID, userID int) error {
	tx, err := cs.Db.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("failed to create transaction")
	}
	defer func() {
		if p := recover(); p != nil || err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("failed to rollback: %v", err)
			}
		}
	}()

	var exists bool
	err = tx.QueryRowContext(context.Background(), "SELECT EXISTS(SELECT 1 FROM chat_members WHERE chat_id=$1 AND user_id=$2)", chatID, userID).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("user %d already exists int chat %d", userID, chatID)
	}

	_, err = tx.ExecContext(context.Background(), "INSERT INTO chat_members(user_id, chat_id) VALUES($1, $2)", userID, chatID)
	if err != nil {
		return fmt.Errorf("failed to execute insert query: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed commit transaction: %w", err)
	}

	return nil
}
