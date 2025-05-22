package models

import "gorm.io/gorm"

type ProductImage struct {
	gorm.Model
	ProductID uint   `json:"product_id" gorm:"not null"`
	ImageURL  string `json:"image_url" gorm:"not null"`
	IsPreview bool   `json:"is_preview" gorm:"default:false"`
	SortOrder int    `json:"sort_order" gorm:"default:0"`
}
