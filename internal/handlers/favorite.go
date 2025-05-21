package handlers

import (
	"alexbirbirdev/go-shop/config"
	"alexbirbirdev/go-shop/internal/models"
	"alexbirbirdev/go-shop/internal/utils"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AddFavorite(c *gin.Context) {
	db := config.DB

	userID, ok := utils.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	var input struct {
		ProductVariantID uint `json:"product_variant_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input - " + err.Error(),
		})
		return
	}

	favorite := models.Favorite{
		UserID:           uint(userID),
		ProductVariantID: input.ProductVariantID,
	}
	var existing models.Favorite
	if err := db.Where("user_id = ? AND product_variant_id = ?", userID, input.ProductVariantID).First(&existing).Error; err == nil {
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
	db := config.DB

	userID, ok := utils.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	limitStr := c.DefaultQuery("limit", "20")
	pageStr := c.DefaultQuery("page", "1")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit
	sortParam := c.DefaultQuery("sort", "created_desc")
	var sortBy string
	switch sortParam {
	case "created_asc":
		sortBy = "created_at ASC"
	case "created_desc":
		fallthrough
	default:
		sortBy = "created_at DESC"
	}

	var favorites []models.Favorite

	if err := db.Offset(offset).Limit(limit).Order(sortBy).Where("user_id = ?", userID).Preload("ProductVariant").Preload("ProductVariant.Product").Find(&favorites).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve favorites",
		})
		return
	}
	type FavoriteResponse struct {
		ID        uint   `json:"product_id"`
		VariantID uint   `json:"variant_id"`
		Name      string `json:"name"`
		// Description string  `json:"description"`
		// Image       string  `json:"image"`
		Price       float64 `json:"price"`
		Stock       int     `json:"stock"`
		VariantName string  `json:"variant_name"`
		IsActive    bool    `json:"is_active"`
	}
	var response []FavoriteResponse
	for _, fav := range favorites {
		p := fav.ProductVariant.Product
		v := fav.ProductVariant
		response = append(response, FavoriteResponse{
			ID:   p.ID,
			Name: p.Name,
			// Description: p.Description,
			// Image:       p.Image,
			Price:       v.Price,
			Stock:       v.Stock,
			VariantName: v.Name,
			VariantID:   v.ID,
			IsActive:    v.IsActive,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"favorites": response,
	})
}

func DeleteFavorite(c *gin.Context) {
	db := config.DB

	userID, ok := utils.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	favoriteID := c.Param("id")

	var favorite models.Favorite
	if err := db.Where("user_id = ? AND product_variant_id = ?", userID, favoriteID).First(&favorite).Error; err != nil {
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
	db := config.DB

	userID, ok := utils.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
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
	db := config.DB

	userID, ok := utils.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	productID := c.Param("id")

	var favorite models.Favorite
	if err := db.Where("user_id = ? AND product_variant_id = ?", userID, productID).First(&favorite).Error; err != nil {
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
