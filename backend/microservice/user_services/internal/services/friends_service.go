package services

import (
	"InterestingChats/backend/user_services/internal/consts"
	"InterestingChats/backend/user_services/internal/models"
	"InterestingChats/backend/user_services/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

func (us *userService) FriendsOperations(r *http.Request, typeRequest int) (*models.Response, int, string, error) {
	var friendRequest models.FriendRequest
	if r.Method != http.MethodGet {
		if err := json.NewDecoder(r.Body).Decode(&friendRequest); err != nil {
			return nil, 400, consts.ErrInternalServer, err
		}
	}

	if _, err := utils.ConvertValue(r.Header.Get("X-User-ID"), &friendRequest.UserID); err != nil {
		return nil, http.StatusUnauthorized, consts.ErrUserUnathorized, fmt.Errorf(consts.InternalTokenInvalid)
	}

	switch typeRequest {
	case consts.FRIENDSHIP_RequestType:
		response, statusCode, clientErr, err := utils.ProxyRequest(consts.POST_Method, consts.DB_AddFriend, friendRequest, http.StatusOK)
		if err != nil {
			return nil, statusCode, clientErr, err
		}

		message, err := utils.MessageConstruct(response, consts.NOTIFICATION_AddFriend)
		if err != nil {
			return nil, http.StatusBadRequest, consts.ErrInternalServer, err
		}

		if err := us.producer.Writer(message, consts.KAFKA_Friendship); err != nil {
			return nil, 400, "", err
		}

		userData, statusCode, clientErr, err := utils.ProxyRequest(consts.GET_Method, fmt.Sprintf(consts.DB_ProfileInfo, friendRequest.UserID), nil, http.StatusOK)
		if err != nil {
			return nil, statusCode, clientErr, err
		}

		if err := us.producer.Writer(models.FriendRequest{UserID: friendRequest.FriendID, FriendID: friendRequest.UserID, Friends: []interface{}{userData.Data}}, consts.KAFKA_PushSubscribers); err != nil {
			return nil, 400, "", err
		}

		return response, http.StatusOK, "", nil
	case consts.FRIENDSHIP_AcceptType:
		response, statusCode, clientErr, err := utils.ProxyRequest(consts.POST_Method, consts.DB_AcceptFriend, friendRequest, http.StatusOK)
		if err != nil {
			return nil, statusCode, clientErr, err
		}

		message, err := utils.MessageConstruct(response, consts.NOTIFICATION_AcceptFriendship)
		if err != nil {
			return nil, http.StatusBadRequest, consts.ErrInternalServer, err
		}

		if err := us.producer.Writer(message, consts.KAFKA_Friendship); err != nil {
			return nil, 400, "", err
		}
		if err := us.producer.Writer(friendRequest, consts.KAFKA_RemoveSubscriberAndAddFriend); err != nil {
			return nil, 400, "", err
		}

		return response, http.StatusOK, "", nil
	case consts.FRIENDSHIP_GetFriendsType:
		response, _, _, err := utils.ProxyRequest(consts.GET_Method, fmt.Sprintf(consts.CACHE_GetFriendList, friendRequest.UserID), nil, http.StatusOK)
		if err != nil || utils.IsEmptyResponseData(response.Data) {
			response, statusCode, clientErr, err := utils.ProxyRequest(consts.GET_Method, fmt.Sprintf(consts.DB_GetFriendList, friendRequest.UserID), nil, http.StatusOK)
			if err != nil {
				return nil, statusCode, clientErr, err
			}

			friendRequest.Friends = response.Data
			if err := us.producer.Writer(friendRequest, consts.KAFKA_PushFriends); err != nil {
				return nil, 400, "error", err
			}

			return response, statusCode, clientErr, err
		}

		return response, http.StatusOK, "", nil
	case consts.FRIENDSHIP_DeleteFriendType:
		response, statusCode, clientErr, err := utils.ProxyRequest(consts.DELETE_Method, consts.DB_DeleteFriend, friendRequest, http.StatusOK)
		if err != nil {
			return nil, statusCode, clientErr, err
		}

		message, err := utils.MessageConstruct(response, consts.NOTIFICATION_DeleteFriend)
		if err != nil {
			return nil, http.StatusBadRequest, consts.ErrInternalServer, err
		}

		if err := us.producer.Writer(message, consts.KAFKA_Friendship); err != nil {
			return nil, 400, "", err
		}
		if err := us.producer.Writer(friendRequest, consts.KAFKA_RemoveFriendAndAddSubscriber); err != nil {
			return nil, 400, "", err
		}

		return response, http.StatusOK, "", nil
	case consts.FRIENDSHIP_DeleteFriendRequestType:
		response, statusCode, clientErr, err := utils.ProxyRequest(consts.DELETE_Method, consts.DB_DeleteFriendRequest, nil, http.StatusOK)
		if err != nil {
			return nil, statusCode, clientErr, err
		}

		return response, http.StatusOK, "", nil
	case consts.FRIENDSHIP_GetSubsType:
		response, _, _, err := utils.ProxyRequest(consts.GET_Method, fmt.Sprintf(consts.CACHE_GetSubs, friendRequest.UserID), nil, http.StatusOK)
		if err != nil || utils.IsEmptyResponseData(response.Data) {
			response, stastatusCode, clientErr, err := utils.ProxyRequest(consts.GET_Method, fmt.Sprintf(consts.DB_GetSubs, friendRequest.UserID), nil, http.StatusOK)
			if err != nil {
				return nil, stastatusCode, clientErr, err
			}

			friendRequest.Friends = response.Data
			if err := us.producer.Writer(friendRequest, consts.KAFKA_PushSubscribers); err != nil {
				return nil, 400, "error", err
			}
			return response, stastatusCode, clientErr, err
		}

		return response, http.StatusOK, "", nil
	default:
		return nil, http.StatusBadRequest, consts.ErrInternalServer, fmt.Errorf("unsupported request type: %d", typeRequest)
	}
}
