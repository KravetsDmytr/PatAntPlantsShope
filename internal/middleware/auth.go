package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type TokenValidator interface {
	ValidateToken(token string) (int, error)
}

func Auth(v TokenValidator) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Потрібна авторизація", "data": nil, "error": "відсутній або невірний токен"})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(header, "Bearer ")
		userID, err := v.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Потрібна авторизація", "data": nil, "error": "невалідний токен"})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
