package rdb

import (
	"InterestingChats/backend/microservice/redis/internal/consts"
	"InterestingChats/backend/microservice/redis/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"
)

// SetSession: save session in Redis
func (rc *RedisClient) SetSession(session *models.Session, ttl int) error {
	sessionData, err := json.Marshal(session)
	if err != nil {
		return err
	}

	sessionKey := fmt.Sprintf("session:%d", session.UserID)
	return rc.client.Set(context.Background(), sessionKey, sessionData, time.Duration(ttl)*time.Second).Err()
}

// GetSession retrieves the session from Redis
func (rc *RedisClient) GetSession(sessionID string) (*models.Session, error) {
	sessionKey := fmt.Sprintf("session:%s", sessionID)
	jsonData, err := rc.client.Get(context.Background(), sessionKey).Result()
	if err != nil {
		return nil, err
	}

	var session models.Session
	if err := json.Unmarshal([]byte(jsonData), &session); err != nil {
		return nil, err
	}
	return &session, nil
}

// UpdateSessionParametr updates certain parameters in the session
func (rc *RedisClient) UpdateSessionParametr(sessionID int, parametr string, value interface{}, ttl int) error {
	sessionKey := fmt.Sprintf("session:%d", sessionID)
	jsonData, err := rc.client.Get(context.Background(), sessionKey).Result()
	if err != nil {
		return err
	}

	var session models.Session
	if err := json.Unmarshal([]byte(jsonData), &session); err != nil {
		return err
	}

	switch parametr {
	case consts.SESSION_Username:
		if username, ok := value.(string); ok && username != session.Username {
			session.Username = username
		}
	case consts.SESSION_Email:
		if email, ok := value.(string); ok && email != session.Email {
			session.Email = email
		}
	case consts.SESSION_Avatar:
		if avatar, ok := value.(string); ok && avatar != *session.Avatar {
			session.Avatar = &avatar
		}
	default:
		return fmt.Errorf(consts.InternalUnknowParametr)
	}

	sessionData, err := json.Marshal(session)
	if err != nil {
		return err
	}
	return rc.client.Set(context.Background(), sessionKey, sessionData, time.Duration(ttl)*time.Second).Err()
}

// GetFriends gets a list of all friends from Redis
func (rc *RedisClient) GetFriends(userID int) ([]*models.Friend, error) {
	pattern := fmt.Sprintf("friends:%d:*", userID)
	friendKeys, err := rc.client.Keys(context.Background(), pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("error fetching friend keys from Redis: %v", err)
	}

	var friends []*models.Friend
	for _, friendKey := range friendKeys {
		friendData, err := rc.client.HGetAll(context.Background(), friendKey).Result()
		if err != nil {
			log.Printf("Error fetching friend data for key %s: %v", friendKey, err)
			continue
		}

		var friend models.Friend
		if id, ok := friendData["id"]; ok {
			friend.ID, _ = strconv.Atoi(id)
		}
		friend.Username = friendData["username"]
		friend.Email = friendData["email"]
		friend.Avatar = friendData["avatar"]

		friends = append(friends, &friend)
	}

	return friends, nil
}

// AddFriends adds friends to Redis, creating a separate key for each friend
func (rc *RedisClient) AddFriends(friends *models.FriendOperation) error {
	for _, friend := range friends.Friend {
		friendKey := fmt.Sprintf("friends:%d:%d", friends.UserID, friend.ID)

		if err := rc.client.HSet(context.Background(), friendKey, map[string]interface{}{
			"id":       friend.ID,
			"username": friend.Username,
			"email":    friend.Email,
			"avatar":   friend.Avatar,
		}).Err(); err != nil {
			return fmt.Errorf("error adding friend to Redis: %v", err)
		}
	}

	return nil
}

// AddSubscribers adds subscribers to Redis, creating a separate key for each subscriber
func (rc *RedisClient) AddSubscribers(subscribers *models.FriendOperation) error {
	for _, friend := range subscribers.Friend {
		subscriberKey := fmt.Sprintf("subscribers:%d:%d", subscribers.UserID, friend.ID)

		if err := rc.client.HSet(context.Background(), subscriberKey, map[string]interface{}{
			"id":       friend.ID,
			"username": friend.Username,
			"email":    friend.Email,
			"avatar":   friend.Avatar,
		}).Err(); err != nil {
			return fmt.Errorf("error adding friend to Redis: %v", err)
		}
	}

	return nil
}

// GetSubscribers извлекает список подписчиков для пользователя
func (rc *RedisClient) GetSubscribers(userID int) ([]*models.Friend, error) {
	pattern := fmt.Sprintf("subscribers:%d:*", userID)
	subsKeys, err := rc.client.Keys(context.Background(), pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("error fetching friend keys from Redis: %v", err)
	}

	var friends []*models.Friend
	for _, friendKey := range subsKeys {
		subData, err := rc.client.HGetAll(context.Background(), friendKey).Result()
		if err != nil {
			log.Printf("Error fetching friend data for key %s: %v", friendKey, err)
			continue
		}

		var friend models.Friend
		if id, ok := subData["id"]; ok {
			friend.ID, _ = strconv.Atoi(id)
		}
		friend.Username = subData["username"]
		friend.Email = subData["email"]
		friend.Avatar = subData["avatar"]

		friends = append(friends, &friend)
	}

	return friends, nil
}

// getListFromSet a helper function to extract and demarshal a list from a Redis set
func (rc *RedisClient) getListFromSet(key string) ([]*models.Friend, error) {
	list, err := rc.client.SMembers(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}

	var friends []*models.Friend
	for _, data := range list {
		var friend models.Friend
		if err := json.Unmarshal([]byte(data), &friend); err != nil {
			return nil, err
		}
		friends = append(friends, &friend)
	}
	return friends, nil
}

// RemoveFriendAndAddSubscriber removes a friend from the friends list and adds it to the subscribers list.
func (rc *RedisClient) RemoveFriendAndAddSubscriber(userID int, friendID int) error {
	friendsKey := fmt.Sprintf("friends:%d:%d", userID, friendID)
	subscribersKey := fmt.Sprintf("subscribers:%d:%d", userID, friendID)

	friendData, err := rc.client.HGetAll(context.Background(), friendsKey).Result()
	if err != nil {
		return err
	}
	if len(friendData) == 0 {
		return fmt.Errorf("friend with ID %d not found", friendID)
	}

	pipe := rc.client.TxPipeline()
	pipe.Del(context.Background(), friendsKey)
	pipe.HSet(context.Background(), subscribersKey, friendData)

	if _, err := pipe.Exec(context.Background()); err != nil {
		return fmt.Errorf("transaction failed: %v", err)
	}

	return nil
}

// RemoveSubscriberAndAddFriend removes a subscriber from the subscriber list and adds them to the friend list.
func (rc *RedisClient) RemoveSubscriberAndAddFriend(userID, subscriberID int) error {
	subscriberKey := fmt.Sprintf("subscribers:%d:%d", userID, subscriberID)
	friendKey := fmt.Sprintf("friends:%d:%d", userID, subscriberID)

	subscriberData, err := rc.client.HGetAll(context.Background(), subscriberKey).Result()
	if err != nil {
		return err
	}
	if len(subscriberData) == 0 {
		return fmt.Errorf("subscriber with ID %d not found", subscriberID)
	}

	pipe := rc.client.TxPipeline()
	pipe.Del(context.Background(), subscriberKey)
	pipe.HSet(context.Background(), friendKey, subscriberData)

	_, err = pipe.Exec(context.Background())
	if err != nil {
		return fmt.Errorf("transaction failed: %v", err)
	}

	return nil
}
