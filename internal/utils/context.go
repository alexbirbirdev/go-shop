package utils

import "github.com/gin-gonic/gin"

func GetUserID(c *gin.Context) (uint, bool) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		return 0, false
	}
	userID, ok := userIDVal.(uint)
	if !ok {
		return 0, false
	}
	return userID, true
}
