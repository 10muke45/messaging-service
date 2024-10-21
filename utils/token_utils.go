package utils

import (
	"errors"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

var secretKey = "your_secret_key"

func GenerateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
	})
	tokenString, err := token.SignedString([]byte(secretKey))
	return tokenString, err
}

// ValidateToken validates the JWT token and returns the username claim if valid
func ValidateToken(tokenString string) (string, bool) {
	// Check if the token starts with "Bearer " and strip it if present
	if strings.HasPrefix(tokenString, "Bearer ") {
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})

	if err != nil || !token.Valid {
		return "", false
	}

	// Extract the claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if username, ok := claims["username"].(string); ok {
			return username, true
		}
	}

	return "", false
}

func GetUsernameFromToken(tokenString string) string {
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	claims, _ := token.Claims.(jwt.MapClaims)
	return claims["username"].(string)
}
