package handlers

import (
	"alexbirbirdev/go-shop/backend/config"
	"alexbirbirdev/go-shop/backend/internal/models"
	"alexbirbirdev/go-shop/backend/internal/utils"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserProfileResponse struct {
	Username   string `json:"username"`
	Name       string `json:"name"`
	LastName   string `json:"last_name"`
	Created_at string `json:"created_at"`
	Role       string `json:"role"`
}

func GetUserProfile(c *gin.Context) {
	db := config.DB
	userID, ok := utils.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}
	var userProfile models.User

	if err := db.First(&userProfile, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve user profile",
		})
		return
	}

	response := UserProfileResponse{
		Username:   userProfile.Email,
		Name:       userProfile.Name,
		LastName:   userProfile.LastName,
		Created_at: userProfile.Model.CreatedAt.Format("2006-01-02 15:04:05"),
		Role:       userProfile.Role,
	}

	c.JSON(http.StatusOK, gin.H{
		"user": response,
	})
}

func UpdateUserProfile(c *gin.Context) {
	db := config.DB

	userID, ok := utils.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized in handler",
		})
		return
	}

	var userProfile models.User

	if err := db.First(&userProfile, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve user profile",
		})
		return
	}

	var input struct {
		Name     string `json:"name"`
		LastName string `json:"last_name"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
		})
		return
	}

	userProfile.Name = input.Name
	userProfile.LastName = input.LastName

	if err := db.Save(&userProfile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update user profile",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User profile updated successfully",
	})

}

// admin
func AdminGetUsers(c *gin.Context) {
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

	var users []models.User

	if err := db.Order(sortBy).Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve users",
		})
		return
	}

	type UserResponse struct {
		ID        uint   `json:"id"`
		Name      string `json:"name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Role      string `json:"role"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	result := make([]UserResponse, len(users))
	for i, user := range users {
		result[i] = UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			LastName:  user.LastName,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt.Format("2006-01-02"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02"),
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"users": result,
	})
}
func AdminGetUser(c *gin.Context) {
	db := config.DB

	id := c.Param("id")

	var user models.User

	if err := db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func AdminUpdateUser(c *gin.Context) {
	db := config.DB

	id := c.Param("id")

	var user models.User

	if err := db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve user",
		})
		return
	}

	var input models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
		})
		return
	}

	if err := db.Model(&user).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update user",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
	})
}
