package handlers

import (
	"alexbirbirdev/go-shop/config"
	"alexbirbirdev/go-shop/internal/models"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateCartItem(c *gin.Context) {
	db := config.InitDB()
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection failed",
		})
		return
	}

	var cartItem models.CartItem

	if err := c.ShouldBindJSON(&cartItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
		})
		return
	}

	cartItem.UserID = 1 // Assuming a logged-in user with ID 1
	var existingCartItem models.CartItem
	if err := db.Where("user_id = ? AND product_id = ? AND product_variant_id = ?", cartItem.UserID, cartItem.ProductID, cartItem.ProductVariantID).First(&existingCartItem).Error; err == nil {
		existingCartItem.Quantity += cartItem.Quantity
		if err := db.Save(&existingCartItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to update cart item",
			})
			return
		}
	} else if err := db.Create(&cartItem).Error; err != nil {
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
	db := config.InitDB()
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection failed",
		})
		return
	}
	var cartItems []models.CartItem
	if err := db.Preload("Product").Preload("ProductVariant").Where("user_id = ?", 1).Find(&cartItems).Error; err != nil {
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

func DeleteCartItem(c *gin.Context) {
	db := config.InitDB()
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection failed",
		})
		return
	}

	id := c.Param("id")
	var cartItem models.CartItem

	if err := db.Where("id = ? AND user_id = ?", id, 1).First(&cartItem).Error; err != nil {
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

	if err := db.Where("id = ? AND user_id = ?", id, 1).Delete(&cartItem).Error; err != nil {
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
	db := config.InitDB()
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection failed",
		})
		return
	}

	user_id := 1 // Assuming a logged-in user with ID 1
	var cart []models.CartItem

	if err := db.Where("user_id = ?", user_id).Find(&cart).Error; err != nil {
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

	if err := db.Where("user_id = ?", user_id).Delete(&cart).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete cart item",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Clear Cart",
	})
}

func UpdateQuantity(c *gin.Context) {
	db := config.InitDB()
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection failed",
		})
		return
	}
	id := c.Param("id")
	var cartItem models.CartItem
	if err := db.Where("id = ? AND user_id = ?", id, 1).First(&cartItem).Error; err != nil {
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
	db := config.InitDB()
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection failed",
		})
		return
	}
	id := c.Param("id")
	var cartItem models.CartItem
	if err := db.Where("id = ? AND user_id = ?", id, 1).First(&cartItem).Error; err != nil {
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
	db := config.InitDB()
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection failed",
		})
		return
	}
	id := c.Param("id")
	var cartItem models.CartItem
	if err := db.Where("id = ? AND user_id = ?", id, 1).First(&cartItem).Error; err != nil {
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
