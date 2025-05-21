package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name            string           `json:"name" gorm:"not null"`
	Description     string           `json:"description" gorm:"not null"`
	Price           float64          `json:"price" gorm:"default:0"`
	Stock           int              `json:"stock" gorm:"default:0"`
	CategoryID      uint             `json:"category_id" gorm:"not null"`
	Images          []ProductImage   `json:"images" gorm:"foreignKey:ProductID"`
	ProductVariants []ProductVariant `json:"product_variants" gorm:"foreignKey:ProductID"`
	IsActive        bool             `json:"is_active" gorm:"default:false"`
}
