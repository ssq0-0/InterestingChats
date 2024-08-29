package rdb

import (
	"InterestingChats/backend/microservice/redis/internal/models"
	"context"
	"encoding/json"
	"log"
)

func (rc *RedisClient) GetTokens(email string) (models.Tokens, error) {
	jsonData, err := rc.client.Get(context.Background(), email).Result()
	if err != nil {
		return models.Tokens{}, err
	}
	var userTokens models.Tokens
	if err = json.Unmarshal([]byte(jsonData), &userTokens); err != nil {
		log.Println("failed decode userTokens data: ", err)
		return models.Tokens{}, err
	}
	return userTokens, nil
}

func (rc *RedisClient) SetToken(userTokens map[string]models.UserTokens) error {
	for email, userToken := range userTokens {
		tokenJSON, err := json.Marshal(userToken.Tokens)
		if err != nil {
			log.Println("Failed to marshal tokens:", err)
			return err
		}
		err = rc.client.Set(context.Background(), email, tokenJSON, 0).Err()
		if err != nil {
			log.Println("Failed to set token in Redis:", err)
			return err
		}
	}
	return nil
}
