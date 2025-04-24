package middleware

import (
	"alexbirbirdev/go-shop/internal/services/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			return
		}

		claims, err := auth.ValidateJWT(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func AuthRequired(c *gin.Context) {
	// Проверяем наличие токена в заголовке или cookies
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		tokenString, _ = c.Cookie("token") // если нет токена в header, ищем его в cookies
	}

	// Если токен не найден, продолжаем обработку запроса
	if tokenString == "" {
		c.Next()
		return
	}

	// Если токен найден, блокируем доступ к эндпоинту
	c.JSON(http.StatusUnauthorized, gin.H{
		"error": "User is already authorized",
	})
	c.Abort()
}
