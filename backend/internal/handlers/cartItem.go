package handlers

import (
	"alexbirbirdev/go-shop/backend/config"
	"alexbirbirdev/go-shop/backend/internal/models"
	"alexbirbirdev/go-shop/backend/internal/utils"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateCartItem(c *gin.Context) {
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
		Quantity         int  `json:"quantity" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
		})
		return
	}

	cartItem := models.CartItem{
		UserID:           uint(userID),
		ProductVariantID: input.ProductVariantID,
		Quantity:         input.Quantity,
	}

	var variant models.ProductVariant
	if err := db.First(&variant, cartItem.ProductVariantID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Product variant not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch product variant",
		})
		return
	}
	if !variant.IsActive {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Product variant is not active",
		})
		return
	}

	var existingCartItem models.CartItem
	if err := db.Where("user_id = ? AND product_variant_id = ?", userID, cartItem.ProductVariantID).First(&existingCartItem).Error; err == nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Cart item not found",
			})
			return
		}
		c.JSON(http.StatusConflict, gin.H{
			"error": "Cart item already exists",
		})
		return
	}
	if err := db.Create(&cartItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create cart item",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Create Cart Item",
	})
}

func GetCartItems(c *gin.Context) {
	db := config.DB
	userID, ok := utils.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}
	var cartItems []models.CartItem
	if err := db.Order("created_at DESC").
		Preload("ProductVariant").
		Preload("ProductVariant.Product").
		Preload("ProductVariant.Product.Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		Where("user_id = ?", userID).Find(&cartItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch cart items",
		})
		return
	}
	var totalPrice float64

	for _, item := range cartItems {
		totalPrice += float64(item.Quantity) * item.ProductVariant.Price
	}

	c.JSON(http.StatusOK, gin.H{
		"cart_items": cartItems,
		"total":      totalPrice,
	})
}
func GetCartNumber(c *gin.Context) {
	db := config.DB
	userID, ok := utils.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	var totalQuantity int64

	if err := db.
		Model(&models.CartItem{}).
		Where("user_id = ?", userID).
		Select("COALESCE(SUM(quantity), 0)").
		Scan(&totalQuantity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при подсчёте"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": totalQuantity})
}

func DeleteCartItem(c *gin.Context) {
	db := config.DB

	userID, ok := utils.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	id := c.Param("id")
	var cartItem models.CartItem

	if err := db.Where("product_variant_id = ? AND user_id = ?", id, userID).First(&cartItem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Cart item not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch cart item",
		})
		return
	}

	if err := db.Where("product_variant_id = ? AND user_id = ?", id, userID).Delete(&cartItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete cart item",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Delete Cart Item",
	})
}

func ClearCart(c *gin.Context) {
	db := config.DB

	userID, ok := utils.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	var cart []models.CartItem

	if err := db.Where("user_id = ?", userID).Find(&cart).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Cart item not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch cart item",
		})
		return
	}

	if err := db.Where("user_id = ?", userID).Delete(&cart).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete cart item",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Cleared",
	})
}

func UpdateQuantity(c *gin.Context) {
	db := config.DB

	userID, ok := utils.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	id := c.Param("id")
	var cartItem models.CartItem
	if err := db.Where("id = ? AND user_id = ?", id, userID).First(&cartItem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Cart item not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch cart item",
		})
		return
	}
	var input struct {
		Quantity int `json:"quantity" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
		})
		return
	}

	var variant models.ProductVariant
	if err := db.Where("id = ?", cartItem.ProductVariantID).First(&variant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Product variant not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch product variant",
		})
		return
	}
	if input.Quantity < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Quantity must be at least 1",
		})
		return
	}
	if input.Quantity > variant.Stock {
		input.Quantity = variant.Stock
	}

	if err := db.Model(&cartItem).Update("Quantity", input.Quantity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update cart item",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Updated quantity",
	})
}

func IncrementCartItem(c *gin.Context) {
	db := config.DB

	userID, ok := utils.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	id := c.Param("id")
	var cartItem models.CartItem
	if err := db.Where("id = ? AND user_id = ?", id, userID).First(&cartItem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Cart item not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch cart item",
		})
		return
	}
	var variant models.ProductVariant
	if err := db.Where("id = ?", cartItem.ProductVariantID).First(&variant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Product variant not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch product variant",
		})
		return
	}
	if cartItem.Quantity >= variant.Stock {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cannot increment quantity beyond stock",
		})
		return
	}
	cartItem.Quantity++

	if err := db.Model(&cartItem).Update("Quantity", cartItem.Quantity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update cart item",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Incremented quantity",
	})
}

func DecrementCartItem(c *gin.Context) {
	db := config.DB

	userID, ok := utils.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	id := c.Param("id")
	var cartItem models.CartItem
	if err := db.Where("id = ? AND user_id = ?", id, userID).First(&cartItem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Cart item not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch cart item",
		})
		return
	}
	if cartItem.Quantity < 2 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cannot decrement quantity below 1",
		})
		return
	}
	cartItem.Quantity--

	if err := db.Model(&cartItem).Update("Quantity", cartItem.Quantity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update cart item",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Incremented quantity",
	})
}

func CheckInCart(c *gin.Context) {
	db := config.DB

	userID, ok := utils.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	response := false
	id := c.Param("id")
	if err := db.Where("id = ? AND user_id = ?", id, userID).First(&models.CartItem{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"in_cart": response,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch cart item",
		})
		return
	}
	response = true
	c.JSON(http.StatusOK, gin.H{
		"in_cart": response,
	})
}
