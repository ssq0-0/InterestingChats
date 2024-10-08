package db

import (
	"InterestingChats/backend/microservice/db/internal/consts"
	"InterestingChats/backend/microservice/db/internal/db_queries"
	"InterestingChats/backend/microservice/db/internal/models"
	"InterestingChats/backend/microservice/db/internal/utils"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/lib/pq"
)

// GetChatInfo retrieves detailed information about a chat based on the chatID and userID.
func (db *dbService) GetChatInfo(chatID, userID int) (*models.Chat, string, error) {
	chat, err := db_queries.SelectChatInfo(db.db, chatID)
	if err != nil {
		return nil, consts.ErrGetChatInfo, err
	}

	isMember, clientErr, err := db_queries.CheckUserChatExists(chat.ID, userID, db.db)
	if err != nil {
		return nil, clientErr, err
	}
	if !isMember {
		return nil, consts.ErrUserNoChatMember, fmt.Errorf(consts.InternalErrUserNoChatMember)
	}

	chat.Members, err = db_queries.SelectChatMembers(chat.ID, db.db)
	if err != nil {
		return nil, consts.ErrScanChatMemberResult, err
	}

	chat.Messages, err = db_queries.SelectChatMessages(chat.ID, db.db)
	if err != nil {
		return nil, consts.ErrScanMessagesResult, err
	}

	return chat, "", nil
}

// GetAllChats retrieves all available chats from the database.
func (db *dbService) GetAllChats() ([]models.Chat, string, error) {
	chatList, clientErr, err := db_queries.SelectAllChats(db.db)
	if err != nil {
		return nil, clientErr, err
	}

	return chatList, "", nil
}

// GetChat retrieves chats that match the specified symbols.
func (db *dbService) GetChat(symbols string) ([]models.Chat, string, error) {
	chatList, clientErr, err := db_queries.SelectAllChatsByQuerySymbol(db.db, symbols)
	if err != nil {
		return nil, clientErr, err
	}

	return chatList, "", nil
}

// GetUserChats fetches all chats associated with a specific user identified by userID.
func (db *dbService) GetUserChats(userID int) ([]models.Chat, string, error) {
	chatList, clientErr, err := db_queries.SelectUserChats(userID, db.db)
	if err != nil {
		return nil, clientErr, err
	}

	return chatList, "", nil
}

// CreateChat creates a new chat based on the given CreateChat model.
func (db *dbService) CreateChat(chat *models.CreateChat) (string, error) {
	tx, err := db.db.BeginTx(context.Background(), nil)
	if err != nil {
		return consts.ErrInternalServerError, err
	}

	defer utils.RollbackOnError(tx, err)

	if chat.ID, err = db_queries.InsertChatInfo(chat.ChatName, chat.Creator, db.db, tx); err != nil || chat.ID == 0 {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return consts.ErrChatAlreadyExists, err
		}
		return consts.ErrChatAlreadyExists, err
	}

	if err := db_queries.InsertChatMembers(chat.ID, tx, chat.Members); err != nil {
		return consts.ErrInternalServerError, err
	}

	if len(chat.Messages) != 0 {
		for _, message := range chat.Messages {
			if err := db_queries.InsertMessage(chat.ID, message.User.ID, message.Body, tx); err != nil {
				return consts.ErrInternalServerError, err
			}
		}
	}

	if len(chat.Hashtags) != 0 {
		for _, tag := range chat.Hashtags {
			if err := db_queries.InsertHashtag(chat.ID, tag.Hashtag, tx); err != nil {
				return consts.ErrInternalServerError, err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return consts.ErrInternalServerError, err
	}
	log.Printf("chat created!")
	return "", nil
}

// DeleteChat removes a chat identified by chatID from the database.
func (db *dbService) DeleteChat(chatID int) (string, error) {
	tx, err := db.db.BeginTx(context.Background(), nil)
	if err != nil {
		return consts.ErrInternalServerError, err
	}

	defer utils.RollbackOnError(tx, err)

	if clientErr, err := db_queries.DeleteChat(chatID, tx); err != nil {
		return clientErr, err
	}

	if err := tx.Commit(); err != nil {
		return consts.ErrInternalServerError, err
	}

	return "", nil
}

// CheckAuthor verifies if the specified userID is the author of the chat identified by chatID.
func (db *dbService) CheckAuthor(userID, chatID int) (string, error) {
	result, err := db_queries.SelectAuthorChat(chatID, db.db)
	if err != nil {
		return consts.ErrNoAuthorFound, err
	}
	if result != userID {
		return consts.ErrUserNoAuthor, fmt.Errorf(consts.InternalErrUserNoAuthor, userID, chatID)
	}

	return "", nil
}

func (db *dbService) DeleteMember(chatID, userID int) (string, error) {
	tx, err := db.db.BeginTx(context.Background(), nil)
	if err != nil {
		return consts.ErrInternalServerError, err
	}

	defer utils.RollbackOnError(tx, err)

	exists, clientErr, err := db_queries.CheckUserChatExists(chatID, userID, db.db)
	if err != nil {
		return clientErr, err
	}
	if !exists {
		return consts.ErrUserNoChatMember, fmt.Errorf(consts.InternalErrUserNoChatMember)
	}

	if err = db_queries.DeleteMemberFromChat(chatID, userID, tx); err != nil {
		return consts.ErrInternalServerError, err
	}

	if err := tx.Commit(); err != nil {
		return consts.ErrInternalServerError, err
	}

	return "", nil
}

// AddMember adds a user to a chat based on chatID and userID.
func (db *dbService) AddMember(chatID, userID int) (string, error) {
	tx, err := db.db.BeginTx(context.Background(), nil)
	if err != nil {
		return consts.ErrInternalServerError, err
	}

	defer utils.RollbackOnError(tx, err)

	if clientErr, err := db_queries.CheckUserExists(db.db, userID); err != nil {
		return clientErr, err
	}

	exists, clientErr, err := db_queries.CheckUserChatExists(chatID, userID, db.db)
	if err != nil {
		return clientErr, err
	}
	if exists {
		return consts.ErrUserChatMember, fmt.Errorf(consts.InternalErrUserAlreadyChatMember, userID, chatID)
	}

	if err = db_queries.InsertMember(userID, chatID, tx); err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return fmt.Sprintf(consts.InternalErrUserAlreadyChatMember, userID, chatID), err
		}
		log.Printf("InsertMember error: %v", err)
		return consts.ErrUserChatMember, fmt.Errorf(consts.InternalErrUserAlreadyChatMember, userID, chatID)
	}

	if err := tx.Commit(); err != nil {
		return consts.ErrInternalServerError, err
	}

	return "", nil
}

// ChangeChatName updates the name of the chat identified by chatID.
func (db *dbService) ChangeChatName(chatID int, chatName string) (string, error) {
	tx, err := db.db.BeginTx(context.Background(), nil)
	if err != nil {
		return consts.ErrInternalServerError, err
	}
	defer utils.RollbackOnError(tx, err)

	if clientErr, err := db_queries.ChangeChatName(chatID, chatName, tx); err != nil {
		return clientErr, err
	}

	if err := tx.Commit(); err != nil {
		return consts.ErrInternalServerError, err
	}
	return "", nil
}

// SaveMessage saves a message to the database using Kafka messaging.
func (db *dbService) SaveMessage(message *kafka.Message) (string, error) {
	var msg *models.Message
	if err := json.Unmarshal(message.Value, &msg); err != nil {
		log.Printf("Failed to decode notification message: %v", err)
		return consts.ErrInvalidValueFormat, err
	}

	tx, err := db.db.BeginTx(context.Background(), nil)
	if err != nil {
		return consts.ErrInternalServerError, err
	}
	defer utils.RollbackOnError(tx, err)

	exists, clientErr, err := db_queries.CheckUserChatExists(msg.ChatID, msg.User.ID, db.db)
	if err != nil {
		return clientErr, err
	}

	if !exists {
		return consts.ErrUserNoChatMember, err
	}

	if err = db_queries.InsertMessage(msg.ChatID, msg.User.ID, msg.Body, tx); err != nil {
		return consts.ErrInternalServerError, err
	}

	if err := tx.Commit(); err != nil {
		return consts.ErrInternalServerError, err
	}

	return "", nil
}

// SetTags sets hashtags for a chat based on the incoming HTTP request.
func (db *dbService) SetTags(r *http.Request) (int, string, error) {
	tx, err := db.db.BeginTx(context.Background(), nil)
	if err != nil {
		return http.StatusBadRequest, consts.ErrInternalServerError, err
	}
	defer utils.RollbackOnError(tx, err)

	var requestInfo *models.MemberRequest
	if err := json.NewDecoder(r.Body).Decode(&requestInfo); err != nil {
		return http.StatusBadRequest, consts.ErrInternalServerError, err
	}

	author, err := db_queries.SelectAuthorChat(requestInfo.ChatID, db.db)
	if err != nil {
		return http.StatusBadRequest, consts.ErrInternalServerError, err
	}
	if author != requestInfo.UserID {
		return http.StatusBadRequest, consts.ErrUserNoAuthor, fmt.Errorf(consts.InternalErrUserNoAuthor, requestInfo.UserID, requestInfo.ChatID)
	}

	for _, tag := range requestInfo.Options {
		if err := db_queries.InsertHashtag(requestInfo.ChatID, tag, tx); err != nil {
			return http.StatusBadRequest, consts.ErrInternalServerError, err
		}
	}

	if err := tx.Commit(); err != nil {
		return http.StatusBadRequest, consts.ErrInternalServerError, err
	}
	return http.StatusOK, "", nil
}

// GetTags retrieves hashtags associated with a chat based on the incoming HTTP request.
func (db *dbService) GetTags(r *http.Request) ([]*models.Hashtag, int, string, error) {
	tags, err := db_queries.SelectTags(r.URL.Query().Get("chatID"), db.db)
	if err != nil {
		return nil, http.StatusBadRequest, consts.ErrMissingParametr, err
	}
	return tags, http.StatusOK, "", nil
}

// DeleteTags removes specified hashtags from a chat based on the incoming HTTP request.
func (db *dbService) DeleteTags(r *http.Request) (int, string, error) {
	tx, err := db.db.BeginTx(context.Background(), nil)
	if err != nil {
		return http.StatusBadRequest, consts.ErrInternalServerError, err
	}
	defer utils.RollbackOnError(tx, err)

	var tags *models.MemberRequest
	if err := json.NewDecoder(r.Body).Decode(&tags); err != nil {
		return http.StatusBadRequest, consts.ErrInvalidRequestFormat, err
	}

	if clientErr, err := db_queries.DeleteTags(tags, tx); err != nil {
		return http.StatusBadRequest, clientErr, err
	}

	if err := tx.Commit(); err != nil {
		return http.StatusBadRequest, consts.ErrInternalServerError, err
	}

	return http.StatusOK, "", nil
}
