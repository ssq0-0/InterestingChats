package utils

import (
	"InterestingChats/backend/user_services/internal/models"
	"errors"
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
	expirationAccessTime := time.Now().Add(5 * time.Minute)
	expirationRefreshTime := time.Now().Add(1 * time.Hour)

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

func ValidateJWT(tokenString string) (int, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return 0, errors.New("token is malformed")
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return 0, errors.New("token is expired")
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return 0, errors.New("token is not valid yet")
			} else {
				return 0, errors.New("could not handle this token: " + err.Error())
			}
		}
		return 0, errors.New("could not handle this token: " + err.Error())
	}

	if !token.Valid {
		return 0, errors.New("invalid token")
	}

	return claims.ID, nil
}

// func RefreshToken(refreshToken string) (string, error) {
// 	claims := &Claims{}

// 	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
// 		return jwtKey, nil
// 	})
// 	if err != nil || !token.Valid {
// 		return "", errors.New("invalid token")
// 	}

// 	newAccessToken, _, err := GenerateJWT(claims.Email)
// 	if err != nil {
// 		return "", err
// 	}

// 	return newAccessToken, nil
// }
