// This code defines the interface and its implementation in the `db` package, which serves as a data access layer
// for the application. It encapsulates methods for working with the database and other services, providing a
// consistent interaction with the various resources.
package db

import (
	"InterestingChats/backend/microservice/db/internal/models"
	"database/sql"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// The `DBService` interface
// This is the key interface that defines a set of methods to handle user data, chats,
// notifications, and tags. Each method provides specific operations such as:
// 1. Create a new user.
// 2. Validating user data.
// 3. Modifying user information, including name and email.
// 4. Chat operations such as creating, deleting, and modifying chat information.
// 5. Managing notifications and interacting with friends.
// This provides a clean and clear architecture for handling data in the application.

// Composition via struct `dbService`.
// The `dbService` structure implements the `DBService` interface and encapsulates the database connection,
// providing concrete implementations of the interface methods. This allows for interaction with the database,
// providing centralised data access control and minimising code duplication.

// The `NewDBService` factory.
// The `NewDBService` function is a factory that creates a new instance of `dbService`, taking as its
// a database connection (of type `*sql.DB`) as an argument. This provides dependency injection,
// which simplifies code testing and maintenance, since different database implementations can be
// easily substituted without changing the rest of the application.
type DBService interface {
	CreateNewUser(user models.User) (*models.User, string, error)
	LoginData(user *models.User) (string, string, error)
	CheckUser(userID int) (bool, string, error)
	GetUserInfo(userID int) (*models.User, string, error)
	ChangeUsername(username string, userID int) (string, error)
	SearchUsers(symbols string) ([]*models.User, string, error)
	ChangeEmail(email string, userID int) (string, error)
	UploadAvatar(req *models.UserFile) (int, string, error)
	FriendshipOperation(r *http.Request, requestType int) (interface{}, int, string, error)

	GetChatInfo(chatID, userID int) (*models.Chat, string, error)
	GetAllChats() ([]models.Chat, string, error)
	GetChat(symbols string) ([]models.Chat, string, error)
	GetUserChats(userID int) ([]models.Chat, string, error)
	CreateChat(chat *models.CreateChat) (string, error)
	DeleteChat(chatID int) (string, error)
	CheckAuthor(userID, chatID int) (string, error)
	DeleteMember(chatID, userID int) (string, error)
	AddMember(chatID, userID int) (string, error)
	ChangeChatName(chatID int, chatName string) (string, error)
	SaveMessage(message *kafka.Message) (string, error)
	SetTags(r *http.Request) (int, string, error)
	GetTags(r *http.Request) ([]*models.Hashtag, int, string, error)
	DeleteTags(r *http.Request) (int, string, error)

	ReadNotifications(r *http.Request) (int, string, error)
	GetNotification(userID string) ([]*models.Notification, int, string, error)
	AddNotification(r *http.Request) (int, string, error)
}

type dbService struct {
	db *sql.DB
}

func NewDBService(db *sql.DB) DBService {
	return &dbService{
		db: db,
	}
}
