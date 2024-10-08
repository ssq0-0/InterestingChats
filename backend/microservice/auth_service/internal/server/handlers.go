package server

import (
	"auth_service/internal/consts"
	"auth_service/internal/models"
	"auth_service/internal/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// Authorization checks the JWT token from the Authorization header,
// validates it, and returns the user ID if the token is valid.
func (s *Server) Authorization(w http.ResponseWriter, r *http.Request) {
	token := strings.ReplaceAll(r.Header.Get("Authorization"), "Bearer ", "")
	if token == "" {
		s.HandleError(w, http.StatusBadRequest, []string{consts.ErrInternalServer}, fmt.Errorf(consts.InternalServerError))
		return
	}

	userID, statusCode, err := utils.ValidateJWT(token, []byte(s.config.JWT.Secret))
	if err != "" || userID == 0 {
		s.HandleError(w, statusCode, []string{consts.ErrInternalServer}, fmt.Errorf(consts.InternalServerError))
		return
	}
	log.Printf("UserID: %d", userID)

	s.SendRespond(w, http.StatusOK, &models.Response{
		Errors: nil,
		Data:   userID,
	})
}

// GenerateTokens generates and returns new JWT access and refresh tokens for a user.
// It decodes the user data from the request body, generates tokens, and sends them back.
func (s *Server) GenerateTokens(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		s.HandleError(w, http.StatusBadRequest, []string{consts.ErrInternalServer}, err)
		return
	}

	tokens, err := utils.GenerateJWT(user, []byte(s.config.JWT.Secret))
	if err != nil {
		s.HandleError(w, http.StatusBadRequest, []string{consts.ErrInternalServer}, err)
		return
	}

	s.SendRespond(w, http.StatusOK, &models.Response{
		Errors: nil,
		Data: models.UserTokens{
			Tokens: *tokens,
		},
	})
}

// RefreshToken validates the provided refresh token and generates a new access token.
// If the refresh token is invalid or expired, it returns an error.
func (s *Server) RefreshToken(w http.ResponseWriter, r *http.Request) {
	oldRefreshToken := r.URL.Query().Get("refreshToken")
	if strings.ReplaceAll(oldRefreshToken, "Bearer ", "") == "" {
		s.HandleError(w, http.StatusBadRequest, []string{consts.ErrUserUnathorized}, fmt.Errorf(consts.InternalTokenError))
		return
	}

	newRefreshToken, clientErr, err := utils.RefreshToken(oldRefreshToken, []byte(s.config.JWT.Secret))
	if err != nil {
		s.HandleError(w, http.StatusBadRequest, []string{clientErr}, err)
		return
	}

	s.SendRespond(w, http.StatusOK, &models.Response{
		Errors: nil,
		Data:   newRefreshToken,
	})
}
