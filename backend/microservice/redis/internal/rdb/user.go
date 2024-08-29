package rdb

import (
	"InterestingChats/backend/microservice/redis/internal/consts"
	"InterestingChats/backend/microservice/redis/internal/models"
	"context"
	"encoding/json"
	"fmt"
)

func (rc *RedisClient) GetTokens(email string) (models.Tokens, error) {
	jsonData, err := rc.client.Get(context.Background(), email).Result()
	if err != nil {
		return models.Tokens{}, fmt.Errorf(consts.InternalFailedRedisRequest, err)
	}

	var userTokens models.Tokens
	if err = json.Unmarshal([]byte(jsonData), &userTokens); err != nil {
		return models.Tokens{}, fmt.Errorf(consts.InternalInvalidValueFormat, err)
	}
	return userTokens, nil
}

func (rc *RedisClient) SetToken(userTokens map[string]models.UserTokens) error {
	for email, userToken := range userTokens {
		tokenJSON, err := json.Marshal(userToken.Tokens)
		if err != nil {
			return fmt.Errorf(consts.InternalInvalidValueFormat, err)
		}
		err = rc.client.Set(context.Background(), email, tokenJSON, 0).Err()
		if err != nil {
			return fmt.Errorf(consts.InternalFailedRedisRequest, err)
		}
	}
	return nil
}

func (rc *RedisClient) DeleteTokens(userID string) error {
	cmd := rc.client.Del(context.Background(), userID)
	if err := cmd.Err(); err != nil {
		return err
	}

	return nil
}
