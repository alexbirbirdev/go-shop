package handlers

import (
	"alexbirbirdev/go-shop/config"
	"alexbirbirdev/go-shop/internal/models"
	"alexbirbirdev/go-shop/internal/services/auth"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserProfileResponse struct {
	Name     string `json:"name"`
	LastName string `json:"last_name"`
}

func GetUserProfile(c *gin.Context) {
	db := config.InitDB()
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection error",
		})
		return
	}
	userID, err := auth.GetUserIDFromCookie(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized - " + err.Error(),
		})
		return
	}
	var userProfile models.User

	if err := db.First(&userProfile, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve user profile",
		})
		return
	}

	response := UserProfileResponse{
		Name:     userProfile.Name,
		LastName: userProfile.LastName,
	}

	c.JSON(http.StatusOK, gin.H{
		"user": response,
	})
}

func UpdateUserProfile(c *gin.Context) {
	db := config.InitDB()
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection error",
		})
		return
	}

	userID, err := auth.GetUserIDFromCookie(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized - " + err.Error(),
		})
		return
	}

	var userProfile models.User

	if err := db.First(&userProfile, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve user profile",
		})
		return
	}

	var input struct {
		Name     string `json:"name"`
		LastName string `json:"last_name"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
		})
		return
	}

	userProfile.Name = input.Name
	userProfile.LastName = input.LastName

	if err := db.Save(&userProfile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update user profile",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User profile updated successfully",
	})

}
