package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateToken(secret string, userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateToken(secret, tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("невірний метод підпису")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, fmt.Errorf("невалідний токен")
	}

	id, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("невалідний user_id")
	}
	return int(id), nil
}
