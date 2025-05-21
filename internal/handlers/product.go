package handlers

import (
	"alexbirbirdev/go-shop/config"
	"alexbirbirdev/go-shop/internal/models"
	"alexbirbirdev/go-shop/internal/utils"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetProducts(c *gin.Context) {
	db := config.DB

	_, limit, offset := getPaginationParams(c)
	sortBy := getSortParam(c)

	// Базовый запрос с фильтрами
	query := db.Model(&models.Product{}).
		Where("is_active = ?", true).
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		Order(sortBy).
		Limit(limit).
		Offset(offset)

	applyFilters(c, query)

	userID, isAuth := utils.GetUserID(c)

	// Сначала получаем все продукты без флагов
	var baseProducts []models.Product
	if err := query.Find(&baseProducts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	// Если пользователь авторизован, проверяем избранное и корзину
	if isAuth {
		// Получаем ID всех продуктов на странице
		var productIDs []uint
		for _, p := range baseProducts {
			productIDs = append(productIDs, p.ID)
		}

		// Получаем продукты в избранном
		var favoriteProductIDs []uint
		if len(productIDs) > 0 {
			db.Model(&models.Favorite{}).
				Select("DISTINCT pv.product_id").
				Joins("JOIN product_variants pv ON pv.id = favorites.product_variant_id AND pv.is_active = true").
				Where("favorites.user_id = ? AND pv.product_id IN (?)", userID, productIDs).
				Pluck("pv.product_id", &favoriteProductIDs)
		}

		// Получаем продукты в корзине
		var cartProductIDs []uint
		if len(productIDs) > 0 {
			db.Model(&models.CartItem{}).
				Select("DISTINCT pv.product_id").
				Joins("JOIN product_variants pv ON pv.id = cart_items.product_variant_id AND pv.is_active = true").
				Where("cart_items.user_id = ? AND pv.product_id IN (?)", userID, productIDs).
				Pluck("pv.product_id", &cartProductIDs)
		}

		// Создаем мапы для быстрой проверки
		favoriteMap := make(map[uint]bool)
		for _, id := range favoriteProductIDs {
			favoriteMap[id] = true
		}

		cartMap := make(map[uint]bool)
		for _, id := range cartProductIDs {
			cartMap[id] = true
		}

		// Формируем ответ с флагами
		type ProductResponse struct {
			ID          uint                  `json:"id"`
			Name        string                `json:"name"`
			Description string                `json:"description"`
			Price       float64               `json:"price"`
			Images      []models.ProductImage `json:"images" gorm:"foreignKey:ProductID"`
			CategoryID  uint                  `json:"category_id"`
			Stock       int                   `json:"stock"`
			IsFavorite  *bool                 `json:"is_favorite,omitempty"`
			InCart      *bool                 `json:"in_cart,omitempty"`
		}

		result := make([]ProductResponse, len(baseProducts))
		for i, product := range baseProducts {
			resp := ProductResponse{
				ID:          product.ID,
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
				Images:      product.Images,
				CategoryID:  product.CategoryID,
				Stock:       product.Stock,
			}

			isFav := favoriteMap[product.ID]
			inCart := cartMap[product.ID]

			resp.IsFavorite = &isFav
			resp.InCart = &inCart

			result[i] = resp
		}

		c.JSON(http.StatusOK, gin.H{
			"products": result,
		})
		return
	}

	// Для неавторизованных пользователей
	type ProductResponse struct {
		ID          uint                  `json:"id"`
		Name        string                `json:"name"`
		Description string                `json:"description"`
		Price       float64               `json:"price"`
		Images      []models.ProductImage `json:"images" gorm:"foreignKey:ProductID"`
		CategoryID  uint                  `json:"category_id"`
		Stock       int                   `json:"stock"`
	}

	result := make([]ProductResponse, len(baseProducts))
	for i, product := range baseProducts {
		result[i] = ProductResponse{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Images:      product.Images,
			CategoryID:  product.CategoryID,
			Stock:       product.Stock,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"products": result,
	})
}

func getPaginationParams(c *gin.Context) (page, limit, offset int) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err = strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}
	offset = (page - 1) * limit
	return
}
func getSortParam(c *gin.Context) string {
	sortParam := c.DefaultQuery("sort", "created_desc")

	switch sortParam {
	case "price_desc":
		return "price DESC"
	case "price_asc":
		return "price ASC"
	case "name_desc":
		return "name DESC"
	case "name_asc":
		return "name ASC"
	case "created_asc":
		return "created_at ASC"
	default:
		return "created_at DESC"
	}
}
func applyFilters(c *gin.Context, query *gorm.DB) {
	if categoryStr := c.Query("category"); categoryStr != "" {
		if category, err := strconv.Atoi(categoryStr); err == nil {
			query.Where("category_id = ?", category)
		}
	}

	if priceMinStr := c.Query("price_min"); priceMinStr != "" {
		if priceMin, err := strconv.ParseFloat(priceMinStr, 64); err == nil {
			query.Where("price >= ?", priceMin)
		}
	}

	if priceMaxStr := c.Query("price_max"); priceMaxStr != "" {
		if priceMax, err := strconv.ParseFloat(priceMaxStr, 64); err == nil {
			query.Where("price <= ?", priceMax)
		}
	}

	if search := c.Query("search"); search != "" {
		query.Where("name ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
}

func GetProduct(c *gin.Context) {
	id := c.Param("id")
	db := config.DB

	// Получаем продукт с вариантами
	var product models.Product
	if err := db.Preload("ProductVariants", "is_active = ?", true).First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product"})
		return
	}

	userID, isAuth := utils.GetUserID(c)

	type VariantResponse struct {
		ID     uint    `json:"id"`
		Name   string  `json:"name"`
		Stock  int     `json:"stock"`
		Price  float64 `json:"price"`
		IsFav  *bool   `json:"is_fav,omitempty"`
		InCart *bool   `json:"in_cart,omitempty"`
	}

	type ProductResponse struct {
		ID          uint    `json:"id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		// Image           string            `json:"image"`
		CategoryID      uint              `json:"category_id"`
		Stock           int               `json:"stock"`
		ProductVariants []VariantResponse `json:"product_variants"`
	}

	// Собираем варианты
	variants := make([]VariantResponse, len(product.ProductVariants))

	if isAuth {
		// Получаем ID всех вариантов продукта
		var variantIDs []uint
		for _, v := range product.ProductVariants {
			variantIDs = append(variantIDs, v.ID)
		}

		// Получаем варианты в избранном
		var favoriteVariantIDs []uint
		if len(variantIDs) > 0 {
			db.Model(&models.Favorite{}).
				Where("user_id = ? AND product_variant_id IN (?)", userID, variantIDs).
				Pluck("product_variant_id", &favoriteVariantIDs)
		}

		// Получаем варианты в корзине
		var cartVariantIDs []uint
		if len(variantIDs) > 0 {
			db.Model(&models.CartItem{}).
				Where("user_id = ? AND product_variant_id IN (?)", userID, variantIDs).
				Pluck("product_variant_id", &cartVariantIDs)
		}

		// Создаем мапы для быстрой проверки
		favoriteMap := make(map[uint]bool)
		for _, id := range favoriteVariantIDs {
			favoriteMap[id] = true
		}

		cartMap := make(map[uint]bool)
		for _, id := range cartVariantIDs {
			cartMap[id] = true
		}

		// Заполняем варианты с флагами
		for i, v := range product.ProductVariants {
			isFav := favoriteMap[v.ID]
			inCart := cartMap[v.ID]

			variants[i] = VariantResponse{
				ID:     v.ID,
				Name:   v.Name,
				Stock:  v.Stock,
				Price:  v.Price,
				IsFav:  &isFav,
				InCart: &inCart,
			}
		}
	} else {
		// Для неавторизованных - без флагов
		for i, v := range product.ProductVariants {
			variants[i] = VariantResponse{
				ID:    v.ID,
				Name:  v.Name,
				Stock: v.Stock,
				Price: v.Price,
			}
		}
	}

	result := ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		// Image:           product.Image,
		CategoryID:      product.CategoryID,
		Stock:           product.Stock,
		ProductVariants: variants,
	}

	c.JSON(http.StatusOK, gin.H{
		"product": result,
	})
}

// admin
func CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
		})
		return
	}

	db := config.DB
	if err := db.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create product",
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Product created",
	})
}

func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	db := config.DB

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
	db := config.DB

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

func AdminGetProducts(c *gin.Context) {
	db := config.DB

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

	search := c.Query("search")
	if search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
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

func AdminGetProduct(c *gin.Context) {
	id := c.Param("id")
	db := config.DB

	var product models.Product
	if err := db.Preload("ProductVariants").First(&product, id).Error; err != nil {
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
