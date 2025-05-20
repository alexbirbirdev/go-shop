package handlers

import (
	"alexbirbirdev/go-shop/config"
	"alexbirbirdev/go-shop/internal/models"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetCategories(c *gin.Context) {
	db := config.InitDB()
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection error",
		})
		return
	}

	var categories []models.Category
	if err := db.Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch categories",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"categories": categories,
	})
}

func GetCategory(c *gin.Context) {
	db := config.InitDB()
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection error",
		})
		return
	}

	id := c.Param("id")
	var category models.Category

	if err := db.First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Category not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch category",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"category": category,
	})
}

// admin
func CreateCategory(c *gin.Context) {
	db := config.InitDB()
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection error",
		})
		return
	}

	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
		})
		return
	}

	if err := db.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create category",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Category created",
	})
}

func DeleteCategory(c *gin.Context) {
	db := config.InitDB()
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection error",
		})
		return
	}
	id := c.Param("id")
	var category models.Category

	if err := db.First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Category not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch category",
		})
		return
	}

	if err := db.Delete(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Category not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Category deleted",
	})
}

func UpdateCategory(c *gin.Context) {
	db := config.InitDB()
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection error",
		})
		return
	}

	id := c.Param("id")

	var category models.Category

	if err := db.First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Category not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch category",
		})
		return
	}

	var input models.Category
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
		})
		return
	}

	if err := db.Model(&category).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update category",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Category updated",
	})
}
