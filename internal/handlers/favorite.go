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

func AddFavorite(c *gin.Context) {
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

	var input struct {
		ProductID uint `json:"product_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input - " + err.Error(),
		})
		return
	}

	favorite := models.Favorite{
		UserID:    uint(userID),
		ProductID: input.ProductID,
	}
	var existing models.Favorite
	if err := db.Where("user_id = ? AND product_id = ?", userID, input.ProductID).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Favorite already exists",
		})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error - " + err.Error(),
		})
		return
	}
	if err := db.Create(&favorite).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to add favorite",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Favorite added successfully",
	})
}

func GetFavorites(c *gin.Context) {
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

	var favorites []models.Favorite

	if err := db.Preload("Product").Where("user_id = ?", userID).Find(&favorites).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve favorites",
		})
		return
	}

	var products []models.Product
	for _, fav := range favorites {
		products = append(products, fav.Product)
	}

	c.JSON(http.StatusOK, gin.H{
		"favorites": products,
	})
}

func DeleteFavorite(c *gin.Context) {
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

	favoriteID := c.Param("id")

	var favorite models.Favorite
	if err := db.Where("user_id = ? AND product_id = ?", userID, favoriteID).First(&favorite).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Favorite not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error - " + err.Error(),
		})
		return
	}

	if err := db.Delete(&favorite).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete favorite",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Favorite deleted successfully",
	})
}

func ClearFavorites(c *gin.Context) {
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
	var favorites []models.Favorite

	if err := db.Where("user_id = ?", userID).Delete(&favorites).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve favorites",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Favorites cleared successfully",
	})
}

func CheckFavorite(c *gin.Context) {
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

	productID := c.Param("id")

	var favorite models.Favorite
	if err := db.Where("user_id = ? AND product_id = ?", userID, productID).First(&favorite).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, gin.H{
				"exists": false,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error - " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"exists": true,
	})
}
