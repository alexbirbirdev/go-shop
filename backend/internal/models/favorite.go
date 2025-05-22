package models

import "gorm.io/gorm"

type Favorite struct {
	gorm.Model
	UserID           uint           `json:"user_id" gorm:"not null"`
	ProductVariantID uint           `json:"product_variant_id" gorm:"not null"`
	ProductVariant   ProductVariant `json:"product_variant" gorm:"foreignKey:ProductVariantID"`
}
