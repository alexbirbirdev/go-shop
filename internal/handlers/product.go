package handlers

import (
	"alexbirbirdev/go-shop/config"
	"alexbirbirdev/go-shop/internal/models"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetProducts(c *gin.Context) {
	db := config.InitDB()
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection error",
		})
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	sortParam := c.DefaultQuery("sort", "created_desc")

	var sortBy string
	switch sortParam {

	case "price_desc":
		sortBy = "price DESC"
	case "price_asc":
		sortBy = "price ASC"
	case "name_desc":
		sortBy = "name DESC"
	case "name_asc":
		sortBy = "name ASC"
	case "created_asc":
		sortBy = "created_at ASC"
	case "created_desc":
		fallthrough
	default:
		sortBy = "created_at DESC"
	}

	categoryStr := c.Query("category")
	priceMinStr := c.Query("price_min")
	priceMaxStr := c.Query("price_max")

	query := db.Model(&models.Product{})

	if categoryStr != "" {
		category, err := strconv.Atoi(categoryStr)
		if err == nil {
			query = query.Where("category_id = ?", category)
		}
	}
	if priceMinStr != "" {
		priceMin, err := strconv.ParseFloat(priceMinStr, 64)
		if err == nil {
			query = query.Where("price >= ?", priceMin)
		}
	}
	if priceMaxStr != "" {
		priceMax, err := strconv.ParseFloat(priceMaxStr, 64)
		if err == nil {
			query = query.Where("price <= ?", priceMax)
		}
	}

	query = query.Order(sortBy).Limit(limit).Offset(offset)

	var products []models.Product
	if err := query.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch products",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
	})
}

func GetProduct(c *gin.Context) {
	id := c.Param("id")
	db := config.InitDB()
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection error",
		})
		return
	}

	var product models.Product
	if err := db.First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Product not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch product",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"product": product,
	})
}

func CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
		})
		return
	}

	db := config.InitDB()
	if err := db.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create product",
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Product created",
		"product": product,
	})
}

func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	db := config.InitDB()
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection error",
		})
		return
	}

	var product models.Product
	if err := db.First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Product not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch product",
		})
		return
	}

	var input models.Product
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
		})
		return
	}

	if err := db.Model(&product).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update product",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Product updated",
	})
}

func DeleteProduct(c *gin.Context) {
	db := config.InitDB()

	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection error",
		})
	}

	id := c.Param("id")

	var product models.Product

	if err := db.First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Product not found",
			})
			return

		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch product",
		})
		return
	}
	if err := db.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete product",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Product deleted",
	})
}
