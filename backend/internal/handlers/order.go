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

func CreateOrder(c *gin.Context) {
	db := config.DB

	userID, ok := utils.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	var cartItems []models.CartItem

	if err := db.Preload("ProductVariant").
		Preload("ProductVariant.Product").
		Where("user_id = ?", userID).
		Find(&cartItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch cart items",
		})
		return
	}

	if len(cartItems) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cart is empty",
		})
		return
	}

	var total float64

	var orderItems []models.OrderItem

	for _, item := range cartItems {
		price := item.ProductVariant.Price
		quantity := item.Quantity
		total += price * float64(quantity)

		orderItem := models.OrderItem{
			ProductVariantID: item.ProductVariantID,
			Quantity:         quantity,
			Price:            price,
		}

		orderItems = append(orderItems, orderItem)
	}

	var input struct {
		Address    string `json:"address" binding:"required"`
		AddressCom string `json:"address_com" binding:"required"`
		Phone      string `json:"phone" binding:"required"`
		FirstName  string `json:"first_name" binding:"required"`
		LastName   string `json:"last_name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
		})
		return
	}

	order := models.Order{
		UserID:     uint(userID),
		TotalPrice: total,
		Status:     "pending",
		OrderItems: orderItems,
		Address:    input.Address,
		AddressCom: input.AddressCom,
		Phone:      input.Phone,
		FirstName:  input.FirstName,
		LastName:   input.LastName,
	}

	if err := db.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create order",
		})
		return
	}

	// Обновляем количество на складе
	for _, item := range orderItems {
		var variant models.ProductVariant
		if err := db.First(&variant, item.ProductVariantID).Error; err != nil {
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
		if variant.Stock < item.Quantity {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Not enough stock for product variant",
			})
			return
		}
		variant.Stock -= item.Quantity

		if variant.Stock == 0 {
			variant.IsActive = false
		}
		if err := db.Save(&variant).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to update product variant stock",
			})
			return
		}
	}

	// Очищаем корзину после создания заказа
	if err := db.Where("user_id = ?", userID).Delete(&cartItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to clear cart",
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Order created successfully",
	})
}

func GetOrders(c *gin.Context) {
	db := config.DB

	userID, ok := utils.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	var orders []models.Order
	if err := db.Order("created_at DESC").
		Preload("OrderItems").
		Preload("OrderItems.ProductVariant").
		Preload("OrderItems.ProductVariant.Product").
		Preload("OrderItems.ProductVariant.Product.Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		Where("user_id = ?", userID).
		Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch orders",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
	})
}

func GetOrder(c *gin.Context) {
	db := config.DB

	userID, ok := utils.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	id := c.Param("id")

	var order models.Order

	if err := db.Preload("OrderItems").
		Preload("OrderItems.ProductVariant").
		Preload("OrderItems.ProductVariant.Product").
		Preload("OrderItems.ProductVariant.Product.Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		Where("user_id = ? AND id = ?", userID, id).Find(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Order not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch order",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"order": order,
	})
}

// admin
func AdminGetOrders(c *gin.Context) {
	db := config.DB

	var orders []models.Order
	if err := db.Order("created_at DESC").
		Preload("OrderItems").
		Preload("OrderItems.ProductVariant").
		Preload("OrderItems.ProductVariant.Product").
		Preload("OrderItems.ProductVariant.Product.Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch orders",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
	})
}

func AdminGetOrder(c *gin.Context) {
	db := config.DB

	id := c.Param("id")

	var order models.Order

	if err := db.Preload("OrderItems").Where("id = ?", id).Find(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Order not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch order",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"order": order,
	})
}

func UpdateOrderStatus(c *gin.Context) {
	db := config.DB

	id := c.Param("id")
	var order models.Order

	if err := db.Where("id = ?", id).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Order not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch order",
		})
		return
	}

	// pending | paid | shipped | delivered | cancelled
	var input struct {
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
		})
		return
	}

	order.Status = input.Status

	if err := db.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update order status",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Order status updated successfully",
	})
}
