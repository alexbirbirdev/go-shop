package auth

import (
	"errors"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID uint, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "go-shop",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ValidateJWT(token string) (*Claims, error) {
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !tkn.Valid {
		return nil, jwt.ErrInvalidKey
	}
	return claims, nil
}

func GetUserIDFromCookie(c *gin.Context) (uint, error) {
	// Получаем токен из cookie
	tokenString, err := c.Cookie("token")
	if err != nil {
		return 0, errors.New("authorization token not found in cookies")
	}

	// Парсим токен и извлекаем claims
	claims, err := ValidateJWT(tokenString)
	if err != nil {
		return 0, err
	}

	// Извлекаем userID из claims
	userID := claims.UserID
	return userID, nil
}
