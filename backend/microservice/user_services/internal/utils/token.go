package utils

import (
	"InterestingChats/backend/user_services/internal/consts"
	"InterestingChats/backend/user_services/internal/models"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("secret")

type Claims struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenerateJWT(user models.User) (string, string, error) {
	expirationAccessTime := time.Now().Add(15 * time.Minute)
	expirationRefreshTime := time.Now().Add(72 * time.Hour)

	accessClaims := &Claims{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationAccessTime.Unix(),
		},
	}

	refreshClaims := &Claims{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationRefreshTime.Unix(),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	refreshtoken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshtoken.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func ValidateJWT(tokenString string) (int, int, string) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return 0, http.StatusBadRequest, consts.ErrTokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return 0, http.StatusUnauthorized, consts.ErrTokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return 0, http.StatusUnauthorized, consts.ErrTokenNotValid
			} else {
				return 0, http.StatusInternalServerError, consts.ErrNotHandleForToken
			}
		}
		return 0, http.StatusBadRequest, consts.ErrNotHandleForToken
	}

	if !token.Valid {
		return 0, http.StatusBadRequest, consts.ErrInvalidToken
	}

	return claims.ID, http.StatusOK, ""
}

func RefreshToken(refreshToken string) (string, string, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return "", consts.ErrInvalidToken, fmt.Errorf(consts.InternalTokenInvalid)
	}

	expirationAccessTime := time.Now().Add(15 * time.Minute)
	newAccessClaims := &Claims{
		ID:       claims.ID,
		Email:    claims.Email,
		Username: claims.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationAccessTime.Unix(),
		},
	}

	newAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newAccessClaims)
	accessTokenString, err := newAccessToken.SignedString(jwtKey)
	if err != nil {
		return "", consts.ErrInternalServer, err
	}

	return accessTokenString, "", nil
}
