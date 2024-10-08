package utils

import (
	"auth_service/internal/consts"
	"auth_service/internal/models"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// GenerateJWT generates a new JWT access and refresh token for the given user.
// The access token expires after 15 minutes, while the refresh token expires after 72 hours.
func GenerateJWT(user models.User, jwtKey []byte) (*models.Tokens, error) {
	expirationAccessTime := time.Now().Add(15 * time.Minute)
	expirationRefreshTime := time.Now().Add(72 * time.Hour)

	accessClaims := &models.Claims{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationAccessTime.Unix(),
		},
	}

	refreshClaims := &models.Claims{
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
		return nil, err
	}

	refreshtoken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshtoken.SignedString(jwtKey)
	if err != nil {
		return nil, err
	}

	return &models.Tokens{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}

// ValidateJWT validates a given JWT token using the provided key.
// It returns the user ID from the token, the HTTP status code, and any error message, if the token is invalid or expired.
func ValidateJWT(tokenString string, jwtKey []byte) (int, int, string) {
	claims := &models.Claims{}

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

// RefreshToken generates a new access token from a valid refresh token.
// If the refresh token is invalid or expired, it returns an error.
func RefreshToken(refreshToken string, jwtKey []byte) (string, string, error) {
	claims := &models.Claims{}

	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return "", consts.ErrInvalidToken, fmt.Errorf(consts.InternalTokenInvalid)
	}

	expirationAccessTime := time.Now().Add(15 * time.Minute)
	newAccessClaims := &models.Claims{
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
