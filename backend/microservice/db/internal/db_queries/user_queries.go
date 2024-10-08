package db_queries

import (
	"InterestingChats/backend/microservice/db/internal/consts"
	"InterestingChats/backend/microservice/db/internal/models"
	"context"
	"database/sql"
	"fmt"
	"log"
)

func SelectUserFromDB(email string, db *sql.DB) (string, error) {
	var existingUserID int
	err := db.QueryRowContext(context.Background(), consts.QUERY_SelectUser, email).Scan(&existingUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return consts.ErrFailedVerifyUser, fmt.Errorf(consts.InternalErrSelectInfo, err)
	}
	if existingUserID != 0 {
		return consts.ErrUserEmailAlreadyExists, fmt.Errorf(consts.InternalErrUserAlreadyExists, email)
	}

	return "", nil
}

func InsertUser(user models.User, db *sql.DB) (*models.User, string, error) {
	if err := db.QueryRowContext(context.Background(), consts.QUERY_InsertUser, user.Username, user.Email, user.Password).Scan(&user.ID); err != nil {
		return nil, consts.ErrInternalServerError, fmt.Errorf(consts.InternalErrFailedInsert, err)
	}

	return &user, "", nil
}

func SelectLoginUserInfo(user *models.User, db *sql.DB) (string, string, error) {
	var password string
	if err := db.QueryRowContext(context.Background(), consts.QUERY_SelectUserData, user.Email).Scan(&user.ID, &user.Username, &password); err != nil {
		return "", consts.ErrUserNotFound, fmt.Errorf(consts.InternalErrUserNotFoud)
	}

	if len(password) == 0 {
		return "", consts.ErrInvalidRequestFormat, fmt.Errorf(consts.InternalErrPasswordNotFound)
	}

	return password, "", nil
}

func ChechUserUxistsById(userID int, db *sql.DB) (string, error) {
	var exists bool
	if err := db.QueryRowContext(context.Background(), consts.QUERY_SelectUserExists, userID).Scan(&exists); err != nil {
		log.Printf("error transaction: %v", err)
		return consts.ErrInternalServerError, err
	}
	if !exists {
		return consts.ErrUserNotFound, fmt.Errorf(consts.InternalErrUserNotFoud)
	}
	return "", nil
}

func SelectAllUserInfo(userID int, db *sql.DB) (*models.User, error) {
	var user models.User
	if err := db.QueryRowContext(context.Background(), consts.QUERY_SelectProfileUserInfo, userID).Scan(&user.Username, &user.Email, &user.Avatar); err != nil {
		return nil, err
	}
	user.ID = userID
	return &user, nil
}

func ChangeUsername(username string, userID int, db *sql.DB) (string, error) {
	var scanResult string
	if err := db.QueryRowContext(context.Background(), consts.QUERY_SelectUsername, userID).Scan(&scanResult); err != nil {
		return consts.ErrInternalServerError, err
	}

	if scanResult == username {
		return consts.ErrUsernameAlredyYouExists, fmt.Errorf(consts.InternalErrChangedUserInfo)
	}

	if _, err := db.ExecContext(context.Background(), consts.QUERY_ChangeUsername, username, userID); err != nil {
		return "", err
	}

	return "", nil
}

func SelectUserBySymbols(symbols string, db *sql.DB) ([]*models.User, string, error) {
	rows, err := db.Query(consts.QUERY_SelectIndexUser, symbols)
	if err != nil {
		return nil, consts.ErrInternalServerError, fmt.Errorf(consts.InternalErrFailedRequest, err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email); err != nil {
			return nil, consts.ErrInternalServerError, err
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, consts.ErrInternalServerError, err
	}

	return users, "", nil
}

func ChangeEmail(email string, userID int, db *sql.DB) (string, error) {
	var existingUserID int
	if err := db.QueryRowContext(context.Background(), consts.QUERY_CheckSingleEmail, email, userID).Scan(&existingUserID); err != nil {
		if err != sql.ErrNoRows {
			return "", err
		}
	} else {
		return consts.ErrEmailAlredyExists, fmt.Errorf(consts.InternalErrChangedUserInfo)
	}

	var scanResult string
	if err := db.QueryRowContext(context.Background(), consts.QUERY_SelectEmail, userID).Scan(&scanResult); err != nil {
		return consts.ErrMissingParametr, err
	}
	if scanResult == email {
		return consts.ErrEmailAlredyYouExists, fmt.Errorf(consts.InternalErrChangedUserInfo)
	}

	if _, err := db.ExecContext(context.Background(), consts.QUERY_ChangeEmail, email, userID); err != nil {
		return "", err
	}

	return "", nil
}

func InsertAvatar(staticURL string, userID int, tx *sql.Tx) (string, error) {
	res, err := tx.ExecContext(context.Background(), consts.QUERY_AddAvatar, staticURL, userID)
	if err != nil {
		return consts.ErrInternalServerError, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return consts.ErrInternalServerError, err
	}
	if rowsAffected == 0 {
		return consts.ErrInternalServerError, nil
	}

	return "", nil
}

func InsertFriendshipRequest(requestID, receiverID int, tx *sql.Tx) (*models.NotificationRequest, string, error) {
	var exists bool
	if err := tx.QueryRowContext(context.Background(), consts.QUERY_SelectFromFriendships, requestID, receiverID).Scan(&exists); err != nil {
		log.Printf("Error executing QUERY_SelectFromFriendships: %v", err)
		return nil, consts.ErrInternalServerError, fmt.Errorf("ошибка проверки существования дружбы: %v", err)
	}
	if exists {
		return nil, consts.ErrAlreadyFriendship, fmt.Errorf("друзья уже: %d и %d", requestID, receiverID)
	}

	var requestExists bool
	if err := tx.QueryRowContext(context.Background(), consts.QUERY_SelectFromFriendshipsRequests, requestID, receiverID).Scan(&requestExists); err != nil {
		log.Printf("Error executing QUERY_SelectFromFriendshipsRequests: %v", err)
		return nil, consts.ErrInternalServerError, fmt.Errorf("ошибка проверки существования запроса на дружбу: %v", err)
	}
	if requestExists {
		return nil, consts.ErrRequestAlreadySent, fmt.Errorf("запрос на дружбу уже отправлен: %d к %d", requestID, receiverID)
	}

	log.Printf("тот, кто добавляет(инициатор):%d кому идет уведомление %d", requestID, receiverID)
	if _, err := tx.ExecContext(context.Background(), consts.QUERY_InsertToFriendshipRequests, requestID, receiverID); err != nil {
		log.Printf("Ошибка вставки в friend_requests: %v", err)
		return nil, consts.ErrInternalServerError, fmt.Errorf("ошибка вставки запроса на дружбу: %v", err)
	}

	var user models.User
	if err := tx.QueryRowContext(context.Background(), consts.QUERY_GetInitator, requestID).Scan(&user.ID, &user.Username); err != nil {
		return nil, consts.ErrInternalServerError, fmt.Errorf("ошибка получения инициатора: %v", err)
	}

	return &models.NotificationRequest{Sender: user, Receiver: receiverID}, "", nil
}

func AcceptFriendshipRequest(requestID, receiverID int, tx *sql.Tx) (*models.NotificationRequest, string, error) {
	log.Printf("req + rec : %d %d", requestID, receiverID)

	var exists bool
	// Проверка существования запроса на дружбу
	if err := tx.QueryRowContext(context.Background(), consts.QUERY_SelectFromFriendshipsRequests, receiverID, requestID).Scan(&exists); err != nil {
		if err == sql.ErrNoRows {
			return nil, consts.ErrFriendshipRequestNotFound, fmt.Errorf("запрос на дружбу не найден для receiverID: %d и requestID: %d", receiverID, requestID)
		}
		return nil, consts.ErrInternalServerError, fmt.Errorf("ошибка проверки существования запроса: %v", err)
	}
	if !exists {
		return nil, consts.ErrFriendshipRequestNotFound, fmt.Errorf("запрос на дружбу не найден для receiverID: %d и requestID: %d", receiverID, requestID)
	}

	// Вставка в таблицу дружбы
	if _, err := tx.ExecContext(context.Background(), consts.QUERY_InsertFriendships, receiverID, requestID); err != nil {
		return nil, consts.ErrInternalServerError, fmt.Errorf("ошибка вставки в friendships: %v", err)
	}

	// Получение информации об инициаторе
	var user models.User
	if err := tx.QueryRowContext(context.Background(), consts.QUERY_GetInitator, requestID).Scan(&user.ID, &user.Username); err != nil {
		return nil, consts.ErrInternalServerError, fmt.Errorf("ошибка получения инициатора: %v", err)
	}

	// Удаление запроса на дружбу
	if _, err := tx.ExecContext(context.Background(), consts.QUERY_DeleteFriendshipRequest, receiverID, requestID); err != nil {
		return nil, consts.ErrInternalServerError, fmt.Errorf("ошибка удаления запроса на дружбу: %v", err)
	}

	return &models.NotificationRequest{Sender: user, Receiver: receiverID}, "", nil
}

func DeleteFriend(remover, removing int, tx *sql.Tx) (*models.NotificationRequest, string, error) {
	var exists bool
	err := tx.QueryRowContext(context.Background(), consts.QUERY_SelectFromFriendships, remover, removing).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, consts.ErrFriendshipNotFound, nil
		}
		return nil, consts.ErrInternalServerError, err
	}

	if !exists {
		return nil, consts.ErrFriendshipNotFound, nil
	}

	log.Printf("тот, кто удалил(инициатор действия) :%d кого удалили(кому придет увед) %d", remover, removing)
	// Удаление дружбы
	if _, err := tx.ExecContext(context.Background(), consts.QUERY_DeleteFriendship, remover, removing); err != nil {
		return nil, consts.ErrInternalServerError, err
	}
	log.Printf("Запрос на вставку: requester_id: %d, receiver_id: %d", removing, remover)
	res, err := tx.ExecContext(context.Background(), consts.QUERY_InsertFriendshipRequest, removing, remover)
	if err != nil {
		log.Printf("Ошибка вставки запроса: %v", err)
		return nil, consts.ErrInternalServerError, fmt.Errorf("oшибка при вставке запроса в friend_requests: %v", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Printf("Ошибка при получении количества затронутых строк: %v", err)
		return nil, consts.ErrInternalServerError, err
	}
	log.Printf("Количество затронутых строк при вставке: %d", rowsAffected)

	// Вставка в friend_requests
	var user models.User
	err = tx.QueryRowContext(context.Background(), consts.QUERY_GetInitator, remover).Scan(&user.ID, &user.Username)
	if err != nil {
		return nil, consts.ErrInternalServerError, err
	}
	// log.Printf("данные пользователя, который удаляет: %v", user)

	return &models.NotificationRequest{Sender: user, Receiver: removing}, "", nil
}

func DeleteFriendRequest(remover, removing int, tx *sql.Tx) (string, error) {
	if _, err := tx.ExecContext(context.Background(), consts.QUERY_DeleteFriendshipRequest, remover, removing); err != nil {
		return consts.ErrInternalServerError, fmt.Errorf(consts.InternalErrToInsert)
	}

	return "", nil
}
func GetFriendList(userID int, db *sql.DB) ([]*models.User, string, error) {
	rows, err := db.Query(consts.QUERY_SelectUsersFriendships, userID)
	if err != nil {
		return nil, consts.ErrInternalServerError, err
	}

	var users []*models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Avatar); err != nil {
			return nil, consts.ErrInternalServerError, err
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, consts.ErrInternalServerError, err
	}

	return users, "", nil
}

func GetSubList(userID int, db *sql.DB) ([]*models.User, string, error) {
	rows, err := db.Query(consts.QUERY_SelectSubscribers, userID)
	if err != nil {
		return nil, consts.ErrInternalServerError, err
	}

	var subscribes []*models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Avatar); err != nil {
			return nil, consts.ErrInternalServerError, err
		}
		subscribes = append(subscribes, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, consts.ErrInternalServerError, err
	}

	return subscribes, "", nil
}
