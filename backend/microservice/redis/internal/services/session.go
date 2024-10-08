package services

import (
	"InterestingChats/backend/microservice/redis/internal/models"
	"InterestingChats/backend/microservice/redis/internal/rdb"
	"encoding/json"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func SetSession(session *kafka.Message, rc *rdb.RedisClient, ttl int) {
	var session_data *models.Session
	if err := json.Unmarshal(session.Value, &session_data); err != nil {
		log.Printf("error: %v", err)
		return
	}

	if err := rc.SetSession(session_data, ttl); err != nil {
		log.Printf("error: %v", err)
		return
	}

}

func AddFriend(friendList *kafka.Message, rc *rdb.RedisClient) {
	var friends *models.FriendOperation
	if err := json.Unmarshal(friendList.Value, &friends); err != nil {
		log.Printf("error: %v", err)
		return
	}

	if err := rc.AddFriends(friends); err != nil {
		log.Printf("error: %v", err)
		return
	}
}

func AddSubscribers(friendList *kafka.Message, rc *rdb.RedisClient) {
	var subscribers *models.FriendOperation
	if err := json.Unmarshal(friendList.Value, &subscribers); err != nil {
		log.Printf("error: %v", err)
		return
	}

	if err := rc.AddSubscribers(subscribers); err != nil {
		log.Printf("error: %v", err)
		return
	}
}

func RemoveFriendAndAddSubs(data *kafka.Message, rc *rdb.RedisClient) {
	var operation models.FriendOperation
	if err := json.Unmarshal(data.Value, &operation); err != nil {
		log.Printf("error unmarshaling data: %v", err)
		return
	}

	if err := rc.RemoveFriendAndAddSubscriber(operation.UserID, operation.FriendID); err != nil {
		log.Printf("error updating friend and subscriber: %v", err)
		return
	}
}

func RemoveSubsAndAddFriend(data *kafka.Message, rc *rdb.RedisClient) {
	var operation models.FriendOperation
	if err := json.Unmarshal(data.Value, &operation); err != nil {
		log.Printf("error unmarshaling data: %v", err)
		return
	}

	if err := rc.RemoveSubscriberAndAddFriend(operation.UserID, operation.FriendID); err != nil {
		log.Printf("error updating subscriber and friend: %v", err)
		return
	}
}

func UpdateSessionData(data *kafka.Message, rc *rdb.RedisClient, ttl int) {
	var changeData models.ChangeUserData
	if err := json.Unmarshal(data.Value, &changeData); err != nil {
		log.Printf("error unmarshaling data: %v", err)
		return
	}

	if err := rc.UpdateSessionParametr(changeData.UserID, changeData.Type, changeData.Data, ttl); err != nil {
		log.Printf("error updating user data: %v", err)
		return
	}
}
