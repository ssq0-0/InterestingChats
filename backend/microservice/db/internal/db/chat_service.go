package db

import (
	"InterestingChats/backend/microservice/db/internal/consts"
	"InterestingChats/backend/microservice/db/internal/db_queries"
	"InterestingChats/backend/microservice/db/internal/models"
	"InterestingChats/backend/microservice/db/internal/utils"
	"context"
	"database/sql"
	"fmt"
	"log"

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

func (cs *ChatService) GetChatInfo(chatID, userID int) (*models.Chat, string, error) {
	chat, err := db_queries.SelectChatInfo(cs.Db, chatID)
	if err != nil {
		return nil, consts.ErrGetChatInfo, err
	}

	isMember, clientErr, err := db_queries.CheckUserChatExists(chat.ID, userID, cs.Db)
	if err != nil {
		return nil, clientErr, err
	}
	if !isMember {
		return nil, consts.ErrUserNoChatMember, fmt.Errorf(consts.InternalErrUserNoChatMember)
	}

	chat.Members, err = db_queries.SelectChatMembers(chat.ID, cs.Db)
	if err != nil {
		return nil, consts.ErrScanChatMemberResult, err
	}

	chat.Messages, err = db_queries.SelectChatMessages(chat.ID, cs.Db)
	if err != nil {
		return nil, consts.ErrScanMessagesResult, err
	}

	return chat, "", nil
}

func (cs *ChatService) GetAllChats() ([]models.Chat, string, error) {
	chatList, clientErr, err := db_queries.SelectAllChats(cs.Db)
	if err != nil {
		return nil, clientErr, err
	}

	return chatList, "", nil
}

func (cs *ChatService) GetChat(symbols string) ([]models.Chat, string, error) {
	chatList, clientErr, err := db_queries.SelectAllChatsByQuerySymbol(cs.Db, symbols)
	if err != nil {
		return nil, clientErr, err
	}

	return chatList, "", nil
}

func (cs *ChatService) GetUserChats(userID int) ([]models.Chat, string, error) {
	chatList, clientErr, err := db_queries.SelectUserChats(userID, cs.Db)
	if err != nil {
		return nil, clientErr, err
	}

	return chatList, "", nil
}

func (cs *ChatService) CreateChat(chat *models.CreateChat) (string, error) {
	tx, err := cs.Db.BeginTx(context.Background(), nil)
	if err != nil {
		return consts.ErrInternalServerError, err
	}

	defer utils.RollbackOnError(tx, err)

	if chat.ID, err = db_queries.InsertChatInfo(chat.ChatName, chat.Creator, cs.Db, tx); err != nil || chat.ID == 0 {
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

	if err := tx.Commit(); err != nil {
		return consts.ErrInternalServerError, err
	}

	return "", nil
}

func (cs *ChatService) DeleteChat(chatID int) (string, error) {
	tx, err := cs.Db.BeginTx(context.Background(), nil)
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

func (cs *ChatService) CheckAuthor(userID, chatID int) (string, error) {
	result, err := db_queries.SelectAuthorChat(chatID, cs.Db)
	if err != nil {
		return consts.ErrNoAuthorFound, err
	}
	if result != userID {
		return consts.ErrUserNoAuthor, fmt.Errorf(consts.InternalErrUserNoAuthor, userID, chatID)
	}

	return "", nil
}

func (cs *ChatService) DeleteMember(chatID, userID int) (string, error) {
	tx, err := cs.Db.BeginTx(context.Background(), nil)
	if err != nil {
		return consts.ErrInternalServerError, err
	}

	defer utils.RollbackOnError(tx, err)

	exists, clientErr, err := db_queries.CheckUserChatExists(chatID, userID, cs.Db)
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

func (cs *ChatService) AddMember(chatID, userID int) (string, error) {
	tx, err := cs.Db.BeginTx(context.Background(), nil)
	if err != nil {
		return consts.ErrInternalServerError, err
	}

	defer utils.RollbackOnError(tx, err)

	if clientErr, err := db_queries.CheckUserExists(cs.Db, userID); err != nil {
		return clientErr, err
	}

	exists, clientErr, err := db_queries.CheckUserChatExists(chatID, userID, cs.Db)
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

func (cs *ChatService) SaveMessage(msg models.Message, chatID int) (string, error) {
	log.Printf("recived request: %+v", msg)
	log.Printf("chatid: %d", chatID)
	tx, err := cs.Db.BeginTx(context.Background(), nil)
	if err != nil {
		return consts.ErrInternalServerError, err
	}
	log.Printf("err: %v", err)
	defer utils.RollbackOnError(tx, err)

	exists, clientErr, err := db_queries.CheckUserChatExists(chatID, msg.User.ID, cs.Db)
	if err != nil {
		log.Printf("err: %v", err)

		return clientErr, err
	}
	log.Printf("exists: %d", exists)

	if !exists {
		return consts.ErrUserNoChatMember, err
	}
	log.Println("here:")

	if err = db_queries.InsertMessage(chatID, msg.User.ID, msg.Body, tx); err != nil {
		log.Printf("err: %v", err)

		return consts.ErrInternalServerError, err
	}
	log.Println("here1")

	if err := tx.Commit(); err != nil {
		log.Printf("err: %v", err)

		return consts.ErrInternalServerError, err
	}
	log.Println("here2:")

	log.Printf("err: %v", err)

	return "", nil
}
