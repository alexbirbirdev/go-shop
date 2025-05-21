package models

import "gorm.io/gorm"

type CartItem struct {
	gorm.Model
	UserID           uint           `json:"user_id"`
	ProductVariantID uint           `json:"product_variant_id" gorm:"not null"`
	ProductVariant   ProductVariant `json:"product_variant" gorm:"foreignKey:ProductVariantID"`
	Quantity         int            `json:"quantity" gorm:"not null"`
}
