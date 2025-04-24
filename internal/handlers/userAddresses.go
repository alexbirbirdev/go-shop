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

func CreateUserAddresses(c *gin.Context) {
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

	var userAddress models.UserAddress
	if err := c.ShouldBindJSON(&userAddress); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
		})
		return
	}
	userAddress.UserID = uint(userID)
	if err := db.Create(&userAddress).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user address",
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "User address created",
	})
}

func GetUserAddress(c *gin.Context) {
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
	id := c.Param("id")

	var userAddresses []models.UserAddress
	if err := db.Where("user_id = ? AND ID = ?", userID, id).Find(&userAddresses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch user addresses",
		})
		return
	}
	c.JSON(http.StatusOK, userAddresses)
}

func GetUserAddresses(c *gin.Context) {
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

	var userAddresses []models.UserAddress
	if err := db.Where("user_id = ?", userID).Find(&userAddresses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch user addresses",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"user_addresses": userAddresses,
	})
}

func DeleteUserAddress(c *gin.Context) {
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
	id := c.Param("id")
	var userAddress models.UserAddress

	if err := db.Where("user_id = ? AND id = ?", userID, id).First(&userAddress).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User address not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete user address",
		})
		return
	}
	if err := db.Where("user_id = ? AND id = ?", userID, id).Delete(&userAddress).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User address not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete user address",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User address deleted",
	})
}
func UpdateUserAddress(c *gin.Context) {
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

	id := c.Param("id")

	var userAddress models.UserAddress
	if err := db.Where("user_id = ? AND id = ?", userID, id).First(&userAddress).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User address not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch user address",
		})
		return
	}
	var input models.UserAddress
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
		})
		return
	}
	if err := db.Model(&userAddress).Where("user_id = ? AND id = ?", userID, id).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update user address",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User address updated",
	})
}
