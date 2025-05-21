package handlers

import (
	"alexbirbirdev/go-shop/config"
	"alexbirbirdev/go-shop/internal/models"
	"errors"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func UploadProductImage(c *gin.Context) {
	db := config.DB
	productID := c.Param("id")

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form data"})
		return
	}

	files := form.File["images"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No files uploaded"})
		return
	}

	var product models.Product
	if err := db.First(&product, productID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	var maxSortOrder int
	if err := db.Model(&models.ProductImage{}).
		Where("product_id = ?", productID).
		Select("COALESCE(MAX(sort_order), 0)").
		Scan(&maxSortOrder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get current sort order"})
		return
	}

	uploadBaseDir := os.Getenv("URL_UPLOAD_PRODUCT_IMAGE")
	if uploadBaseDir == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Upload directory not configured"})
		return
	}

	uploadDir := uploadBaseDir + productID
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	var images []models.ProductImage
	var previewExists bool

	var existingPreviewCount int64
	if err := db.Model(&models.ProductImage{}).
		Where("product_id = ? AND is_preview = ?", productID, true).
		Count(&existingPreviewCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing preview"})
		return
	}

	for i, file := range files {
		ext := filepath.Ext(file.Filename)
		fileName := uuid.New().String() + ext
		filePath := filepath.Join(uploadDir, fileName)

		if err := c.SaveUploadedFile(file, filePath); err != nil {
			continue
		}

		isPreview := i == 0 && existingPreviewCount == 0 && !previewExists

		images = append(images, models.ProductImage{
			ProductID: product.ID,
			ImageURL:  uploadDir + "/" + fileName,
			IsPreview: isPreview,
			SortOrder: maxSortOrder + 1 + i,
		})

		if isPreview {
			previewExists = true
		}
	}

	// Сохраняем в транзакции
	tx := db.Begin()
	if err := tx.Create(&images).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save images to database"})
		return
	}
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message": "Images uploaded successfully",
	})
}

func SetPreviewImage(c *gin.Context) {
	db := config.DB
	productID := c.Param("id")
	imageID := c.Param("image_id")

	tx := db.Begin()

	var oldPreview models.ProductImage
	if err := db.Model(&models.ProductImage{}).
		Where("product_id = ? AND is_preview = ?", productID, true).
		First(&oldPreview).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch old preview image",
		})
		return
	}

	var newPreview models.ProductImage
	if err := db.Model(&models.ProductImage{}).
		Where("product_id = ? AND id = ?", productID, imageID).
		First(&newPreview).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{
			"error": "New preview image not found",
		})
		return
	}

	if oldPreview.ID != 0 {
		if err := tx.Model(&oldPreview).
			Updates(map[string]interface{}{
				"is_preview": false,
				"sort_order": newPreview.SortOrder,
			}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to update old preview image",
			})
			return
		}
	}

	if err := tx.Model(&newPreview).
		Updates(map[string]interface{}{
			"is_preview": true,
			"sort_order": 0,
		}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update new preview image",
		})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message": "Preview image updated",
	})
}

func UpdateImageOrder(c *gin.Context) {
	db := config.DB
	productID := c.Param("id")

	var input struct {
		ImageIDs []uint `json:"image_ids"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input - " + err.Error(),
		})
		return
	}

	tx := db.Begin()
	for order, imageID := range input.ImageIDs {
		if err := tx.Model(&models.ProductImage{}).Where("product_id = ? AND id = ?", productID, imageID).Update("sort_order", order).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to update image order",
			})
			return
		}
	}
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message": "Image order updated successfully",
	})
}

func DeleteProductImage(c *gin.Context) {
	db := config.DB
	productID := c.Param("id")
	imageID := c.Param("image_id")

	if err := db.Model(&models.Product{}).Where("id = ?", productID).First(&models.Product{}).Error; err != nil {
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
	if err := db.Model(&models.ProductImage{}).Where("id = ?", imageID).First(&models.ProductImage{}).Error; err != nil {
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

	if err := db.Where("product_id = ? AND id = ?", productID, imageID).Delete(&models.ProductImage{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete image",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Image deleted successfully",
	})
}

func GetProductImages(c *gin.Context) {
	db := config.DB
	productID := c.Param("id")

	if err := db.Model(&models.Product{}).Where("id = ?", productID).First(&models.Product{}).Error; err != nil {
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

	var images []models.ProductImage

	if err := db.Order("sort_order ASC").Where("product_id = ?", productID).Find(&images).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch images",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"images": images,
	})
}
