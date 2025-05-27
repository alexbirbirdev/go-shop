package handlers

import (
	"alexbirbirdev/go-shop/backend/config"
	"alexbirbirdev/go-shop/backend/internal/models"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetProductVariants(c *gin.Context) {
	db := config.DB

	id := c.Param("id")
	var variants []models.ProductVariant
	if err := db.Where("product_id = ?", id).Find(&variants).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch variants",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"variants": variants,
	})

}

func GetProductVariant(c *gin.Context) {
	db := config.DB

	id := c.Param("id")
	var variant models.ProductVariant

	if err := db.First(&variant, id).Error; err != nil {
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
	c.JSON(http.StatusOK, gin.H{
		"variant": variant,
	})
}

// admin

func CreateProductVariant(c *gin.Context) {
	db := config.DB
	productIDParam := c.Param("id")
	type VariantsInput struct {
		Variants []models.ProductVariant `json:"variants"`
	}
	var input VariantsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат запроса: " + err.Error(),
		})
		return
	}

	var product models.Product
	if err := db.First(&product, productIDParam).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"error": "Товар не найден",
		})
		return
	}
	for i := range input.Variants {
		input.Variants[i].ProductID = product.ID
	}

	if err := db.Create(&input.Variants).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Один из вариантов уже существует для этого товара",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Не удалось сохранить варианты: " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message":  "Варианты успешно созданы",
		"variants": input.Variants,
	})
}

// func CreateProductVariant(c *gin.Context) {
// 	db := config.DB

// 	id := c.Param("id")
// 	var products []models.Product
// 	if err := db.First(&products, id).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Failed to fetch product",
// 		})
// 		return
// 	}

// 	var variant models.ProductVariant

// 	if err := c.ShouldBindJSON(&variant); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Invalid input",
// 		})
// 		return
// 	}

// 	variant.ProductID = products[0].ID

// 	if err := db.Create(&variant).Error; err != nil {
// 		if strings.Contains(err.Error(), "duplicate") {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Variant with this name already exists for the product"})
// 			return
// 		}
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create variant"})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, gin.H{
// 		"message": "Product variant created",
// 	})
// }

func DeleteProductVariant(c *gin.Context) {
	db := config.DB

	id := c.Param("id")
	var variant models.ProductVariant
	if err := db.First(&variant, id).Error; err != nil {
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
	if err := db.Delete(&variant, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete product variant",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Product variant deleted",
	})
}

func UpdateProductVariant(c *gin.Context) {
	db := config.DB

	id := c.Param("id")

	var variant models.ProductVariant
	if err := db.First(&variant, id).Error; err != nil {
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

	var input struct {
		Name  *string  `json:"name,omitempty"`
		Stock *int     `json:"stock,omitempty"`
		Price *float64 `json:"price,omitempty"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Создаем карту для обновления только тех полей, которые пришли
	updates := make(map[string]interface{})

	if input.Name != nil {
		updates["name"] = *input.Name
	}
	if input.Stock != nil {
		updates["stock"] = *input.Stock
		// Также обновляем флаг IsActive
		updates["is_active"] = (*input.Stock > 0)
	}
	if input.Price != nil {
		updates["price"] = *input.Price
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	if err := db.Model(&variant).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update product variant",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product variant updated",
	})
}
