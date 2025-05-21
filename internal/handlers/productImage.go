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
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to parse form data",
		})
		return
	}

	files := form.File["images"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No files uploaded",
		})
		return
	}

	var product models.Product
	if err := db.First(&product, productID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Product not found",
		})
		return
	}

	uploadBaseDir := os.Getenv("URL_UPLOAD_PRODUCT_IMAGE")
	if uploadBaseDir == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve upload directory from environment",
		})
		return
	}
	uploadDir := uploadBaseDir + productID
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create upload directory",
		})
		return
	}

	var images []models.ProductImage
	var previewExists bool
	for i, file := range files {
		ext := filepath.Ext(file.Filename)
		fileName := uuid.New().String() + ext
		filePath := filepath.Join(uploadDir+"/", fileName)

		if err := c.SaveUploadedFile(file, filePath); err != nil {
			continue
		}

		isPreview := i == 0 && !previewExists

		images = append(images, models.ProductImage{
			ProductID: product.ID,
			ImageURL:  uploadDir + "/" + fileName,
			IsPreview: isPreview,
			SortOrder: len(images),
		})

		if isPreview {
			previewExists = true
		}
	}

	if err := db.Create(&images).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save images to database",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Images uploaded successfully",
		"images":  images,
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

	if err := tx.Model(&models.ProductImage{}).
		Where("product_id = ? AND id NOT IN (?)", productID, []uint{newPreview.ID, oldPreview.ID}).
		Update("sort_order", gorm.Expr("sort_order + 1")).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update sort order of other images",
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
