package db_queries

import (
	"InterestingChats/backend/microservice/db/internal/consts"
	"InterestingChats/backend/microservice/db/internal/models"
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/lib/pq"
)

func SelectAllChats(db *sql.DB) ([]models.Chat, string, error) {
	rows, err := db.Query(consts.QUERY_SelectAllChats)
	if err != nil {
		return nil, consts.ErrInternalServerError, fmt.Errorf(consts.InternalErrFailedRequest, err)
	}
	defer rows.Close()

	var chats []models.Chat
	for rows.Next() {
		var chat models.Chat
		if err := rows.Scan(&chat.ID, &chat.ChatName, &chat.Creator); err != nil {
			return nil, consts.ErrInternalServerError, err
		}
		chats = append(chats, chat)
	}

	if err := rows.Err(); err != nil {
		return nil, consts.ErrInternalServerError, err
	}

	return chats, "", nil
}

func SelectAllChatsByQuerySymbol(db *sql.DB, symbols string) ([]models.Chat, string, error) {
	queryString := symbols + ":*"
	rows, err := db.Query(consts.QUERY_SelectIndexChat, queryString)
	if err != nil {
		return nil, consts.ErrInternalServerError, fmt.Errorf(consts.InternalErrFailedRequest, err)
	}
	defer rows.Close()

	var chats []models.Chat
	for rows.Next() {
		var chat models.Chat
		if err := rows.Scan(&chat.ID, &chat.ChatName, &chat.Creator); err != nil {
			return nil, consts.ErrInternalServerError, err
		}
		chats = append(chats, chat)
	}

	if err := rows.Err(); err != nil {
		return nil, consts.ErrInternalServerError, err
	}

	return chats, "", nil
}

func SelectUserChats(userID int, db *sql.DB) ([]models.Chat, string, error) {
	rows, err := db.Query(consts.QUERY_SelectUsersChat, userID)
	if err != nil {
		return nil, consts.ErrInternalServerError, fmt.Errorf(consts.InternalErrFailedRequest, err)
	}
	defer rows.Close()

	var chats []models.Chat
	for rows.Next() {
		var chat models.Chat
		if err := rows.Scan(&chat.ID, &chat.ChatName, &chat.Creator); err != nil {
			return nil, consts.ErrInternalServerError, err
		}
		chats = append(chats, chat)
	}

	if err := rows.Err(); err != nil {
		return nil, consts.ErrInternalServerError, err
	}

	return chats, "", nil
}

func SelectChatInfo(db *sql.DB, chatID int) (*models.Chat, error) {
	chat := &models.Chat{
		ID:      chatID,
		Members: make(map[int]models.User),
	}

	if err := db.QueryRowContext(context.Background(), consts.QUERY_SelectChatInfo, chatID).Scan(&chat.ChatName, &chat.Creator); err != nil {
		return nil, fmt.Errorf(consts.InternalErrSelectInfo, err)
	}
	return chat, nil
}

func CheckUserExists(db *sql.DB, userID int) (string, error) {
	var exists bool
	if err := db.QueryRowContext(context.Background(), consts.QUERY_SelectUserExists, userID).Scan(&exists); err != nil {
		return consts.ErrUserNoExists, fmt.Errorf(consts.InternalErrCheckUser, err)
	}
	if !exists {
		return consts.ErrUserNoExists, fmt.Errorf(consts.InternalErrUserNotFoud)
	}

	return "", nil
}

func SelectChatMembers(chatID int, db *sql.DB) (map[int]models.User, error) {
	memberRows, err := db.QueryContext(context.Background(), consts.QUERY_SelectChatMembers, chatID)
	if err != nil {
		return nil, fmt.Errorf(consts.InternalErrSelectInfo, err)
	}
	defer memberRows.Close()

	chatMembers := make(map[int]models.User)
	for memberRows.Next() {
		var user models.User
		if err := memberRows.Scan(&user.ID, &user.Username, &user.Email, &user.Avatar); err != nil {
			return nil, fmt.Errorf(consts.InternalErrScanResult, err)
		}

		chatMembers[user.ID] = user
	}

	if err = memberRows.Err(); err != nil {
		return nil, fmt.Errorf(consts.InternalErrIterateErrors, err)
	}

	return chatMembers, nil
}

func SelectChatMessages(chatID int, db *sql.DB) ([]models.Message, error) {
	messageRows, err := db.QueryContext(context.Background(), consts.QUERY_SelectChatMessages, chatID)
	if err != nil {
		return nil, fmt.Errorf(consts.InternalErrSelectInfo, err)
	}
	defer messageRows.Close()

	var messages []models.Message
	for messageRows.Next() {
		var msg models.Message
		if err := messageRows.Scan(&msg.ID, &msg.Body, &msg.Time, &msg.User.ID, &msg.User.Username, &msg.User.Email); err != nil {
			return nil, fmt.Errorf(consts.InternalErrScanResult, err)
		}
		messages = append(messages, msg)
	}

	if err := messageRows.Err(); err != nil {
		return nil, fmt.Errorf(consts.InternalErrIterateErrors, err)
	}

	return messages, nil
}

func InsertChatInfo(chatName string, creator int, db *sql.DB, tx *sql.Tx) (int, error) {
	var chatID int
	if err := tx.QueryRowContext(context.Background(), consts.QUERY_InsertChat, chatName, time.Now(), creator).Scan(&chatID); err != nil {
		return 0, fmt.Errorf(consts.InternalErrFailedInsert, err)
	}

	return chatID, nil
}

func InsertChatMembers(chatID int, tx *sql.Tx, members []int) error {
	insertMemberStmt, err := tx.PrepareContext(context.Background(), consts.QUERY_InsertMembers)
	if err != nil {
		return fmt.Errorf(consts.InternalErrPrepareTx, err)
	}
	defer insertMemberStmt.Close()

	for _, userID := range members {
		if _, err := insertMemberStmt.ExecContext(context.Background(), chatID, userID); err != nil {
			return fmt.Errorf(consts.InternalErrFailedInsert, err)
		}
	}

	return nil
}

func SelectAuthorChat(chatID int, db *sql.DB) (int, error) {
	var result int
	if err := db.QueryRowContext(context.Background(), consts.QUERY_SelectAuthor, chatID).Scan(&result); err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf(consts.InternalErrAuthorNoFound, err)
		}
		return 0, fmt.Errorf(consts.InternalErrFailedInsert, err)
	}

	return result, nil
}

func DeleteMemberFromChat(chatID, userID int, tx *sql.Tx) error {
	if _, err := tx.ExecContext(context.Background(), consts.QUERY_DeleteMemberFromChat, chatID, userID); err != nil {
		return fmt.Errorf(consts.InternalErrDelete, err)
	}
	return nil
}

func DeleteChat(chatID int, tx *sql.Tx) (string, error) {
	if _, err := tx.ExecContext(context.Background(), consts.QUERY_DeleteChat, chatID); err != nil {
		log.Printf("failed to exec: %v", err)
		return consts.ErrFailedDeleteChat, fmt.Errorf(consts.InternalErrDelete, err)
	}

	return "", nil
}

func InsertMember(userID, chatID int, tx *sql.Tx) error {
	_, err := tx.ExecContext(context.Background(), consts.QUERY_InsertMembers, chatID, userID)
	if err != nil {
		log.Printf("Error during insert operation: %v", err)
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505":
				return fmt.Errorf("user %d is already a member of chat %d", userID, chatID)
			default:
				return fmt.Errorf("PostgreSQL error: %v", pqErr)
			}
		}
		return fmt.Errorf("failed to insert member: %v", err)
	}
	return nil
}

func InsertMessage(chatID, userID int, body string, tx *sql.Tx) error {
	if sqlRes, err := tx.ExecContext(context.Background(), consts.QUERY_InsertMessage, chatID, userID, body, time.Now()); err != nil {
		log.Println("sql res: ", sqlRes)
		return fmt.Errorf(consts.InternalErrFailedInsert, err)
	}

	return nil
}

func InsertHashtag(chatID int, tag string, tx *sql.Tx) error {
	var hashtagID int
	err := tx.QueryRowContext(context.Background(), consts.QUERY_InsertInHashtag, tag).Scan(&hashtagID)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("ошибка при вставке хештега: %v", err)
	}

	if err == sql.ErrNoRows {
		err = tx.QueryRowContext(context.Background(), "SELECT id FROM hashtags WHERE name = $1", tag).Scan(&hashtagID)
		if err != nil {
			return fmt.Errorf("не удалось найти хештег: %v", err)
		}
	}

	if _, err = tx.ExecContext(context.Background(), consts.QUERY_InsertInChatHashtags, chatID, hashtagID); err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return fmt.Errorf(consts.InternalTagAlreadySet) // Связь уже существует, ошибок нет
		}
		return fmt.Errorf("ошибка при вставке связи хештега с чатом: %v", err)
	}

	return nil
}

func SelectTags(chatID string, db *sql.DB) ([]*models.Hashtag, error) {
	tagRows, err := db.QueryContext(context.Background(), consts.QUERY_SelectTags, chatID)
	if err != nil {
		return nil, err
	}
	defer tagRows.Close()

	var tagList []*models.Hashtag
	for tagRows.Next() {
		var tag models.Hashtag
		if err := tagRows.Scan(&tag.ID, &tag.Hashtag); err != nil {
			return nil, fmt.Errorf(consts.InternalErrScanResult, err)
		}
		tagList = append(tagList, &tag)
	}

	if err := tagRows.Err(); err != nil {
		return nil, fmt.Errorf(consts.InternalErrIterateErrors, err)
	}

	return tagList, nil
}

func DeleteTags(tags *models.MemberRequest, tx *sql.Tx) (string, error) {
	if len(tags.Options) == 0 {
		return consts.ErrTagsNotFound, nil
	}

	if _, err := tx.ExecContext(context.Background(), consts.QUERY_DeleteTags, tags.ChatID, pq.Array(tags.Options)); err != nil {
		return consts.ErrInternalServerError, fmt.Errorf("ошибка при удалении тегов: %w", err)
	}

	return "", nil
}

func CheckUserChatExists(chatID, userID int, db *sql.DB) (bool, string, error) {
	log.Printf("userID & chatID in func: %d %d", userID, chatID)
	var isMember bool
	if err := db.QueryRowContext(context.Background(), consts.QUERY_UserChatExists, chatID, userID).Scan(&isMember); err != nil {
		log.Printf("err: %v", err)
		return false, consts.ErrUserChatExists, fmt.Errorf(consts.InternalErrCheckUser, err)
	}
	log.Println(isMember)
	return isMember, "", nil
}

func ChangeChatName(chatID int, chatName string, tx *sql.Tx) (string, error) {
	if _, err := tx.ExecContext(context.Background(), consts.QUERY_ChangeChatName, chatName, chatID); err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			return consts.ErrChatNameAlreadyExists, fmt.Errorf(consts.InternalChatNameExists)
		}
		return consts.ErrInternalServerError, err
	}

	return "", nil
}
