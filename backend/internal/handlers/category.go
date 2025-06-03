package handlers

import (
	"alexbirbirdev/go-shop/backend/config"
	"alexbirbirdev/go-shop/backend/internal/models"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetCategories(c *gin.Context) {
	db := config.DB

	parentIDStr := c.Query("parent_id")

	var categories []models.Category

	if parentIDStr == "" {
		if err := db.Preload("Children").Where("parent_id IS NULL").Find(&categories).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch categories",
			})
			return
		}
	} else {
		parentID, err := strconv.ParseUint(parentIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid parent_id",
			})
			return
		}
		if err := db.Preload("Children").Where("id = ?", parentID).Find(&categories).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch categories",
			})
			return
		}
	}
	type CategoryResponse struct {
		ID       uint   `json:"id"`
		Name     string `json:"name"`
		ParentID *uint  `json: "parent_id"`
		Children []struct {
			ID   uint   `json:"id"`
			Name string `json:"name"`
		}
	}
	var response []CategoryResponse
	for _, c := range categories {
		response = append(response, CategoryResponse{
			ID:       c.ID,
			Name:     c.Name,
			ParentID: c.ParentID,
			Children: []struct {
				ID   uint   `json:"id"`
				Name string `json:"name"`
			}{},
		})
		for _, child := range c.Children {
			response[len(response)-1].Children = append(response[len(response)-1].Children, struct {
				ID   uint   `json:"id"`
				Name string `json:"name"`
			}{
				ID:   child.ID,
				Name: child.Name,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"categories": response,
	})
}

func GetCategory(c *gin.Context) {
	db := config.DB

	id := c.Param("id")
	var category models.Category

	if err := db.Preload("Children").First(&category, id).Error; err != nil {
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

	type CategoryResponse struct {
		ID       uint   `json:"id"`
		Name     string `json:"name"`
		ParentID *uint  `json:"parent_id"`
		// Children []struct {
		// 	ID   uint   `json:"id"`
		// 	Name string `json:"name"`
		// }
	}
	var response CategoryResponse
	response.ID = category.ID
	response.Name = category.Name
	response.ParentID = category.ParentID
	// for _, child := range category.Children {
	// 	response.Children = append(response.Children, struct {
	// 		ID   uint   `json:"id"`
	// 		Name string `json:"name"`
	// 	}{
	// 		ID:   child.ID,
	// 		Name: child.Name,
	// 	})
	// }

	c.JSON(http.StatusOK, gin.H{
		"category": response,
	})
}

// admin
func CreateCategory(c *gin.Context) {
	db := config.DB

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
	db := config.DB
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
	db := config.DB

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
	// Check if the category name already exists
	// if input.Name == "" {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error": "Category name is required",
	// 	})
	// 	return
	// }

	if input.Name != "" {
		var existingCategory models.Category
		if err := db.Where("name = ?", input.Name).First(&existingCategory).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Category with this name already exists",
			})
			return
		}
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

func AdminGetCategories(c *gin.Context) {
	db := config.DB

	parentIDStr := c.Query("parent_id")

	var categories []models.Category

	if parentIDStr == "" {
		if err := db.Preload("Children").Where("parent_id IS NULL").Find(&categories).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch categories",
			})
			return
		}
	} else {
		parentID, err := strconv.ParseUint(parentIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid parent_id",
			})
			return
		}
		if err := db.Preload("Children").Where("id = ?", parentID).Find(&categories).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch categories",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"categories": categories,
	})
}

func AdminGetCategory(c *gin.Context) {
	db := config.DB

	id := c.Param("id")
	var category models.Category

	if err := db.Preload("Children").First(&category, id).Error; err != nil {
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

func ShowAllCategories(c *gin.Context) {
	db := config.DB

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	page, err := strconv.Atoi(pageStr)
	if page < 1 || err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if limit < 1 || err != nil {
		limit = 50
	}
	offset := (page - 1) * limit

	sortParam := c.DefaultQuery("sort", "name_ASC")
	var sortBy string
	switch sortParam {
	case "created_asc":
		sortBy = "created_at ASC"
	case "created_desc":
		sortBy = "created_at DESC"
	case "name_desc":
		sortBy = "name DESC"
	case "name_asc":
		fallthrough
	default:
		sortBy = "name ASC"
	}

	var categories []models.Category
	if err := db.Order(sortBy).Offset(offset).Limit(limit).Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch categories",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"categories": categories,
	})
}
