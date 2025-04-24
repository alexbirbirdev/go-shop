package models

import "gorm.io/gorm"

type CartItem struct {
	gorm.Model
	UserID           uint           `json:"user_id"`
	ProductID        uint           `json:"product_id" gorm:"not null"`
	ProductVariantID uint           `json:"product_variant_id" gorm:"not null"`
	Quantity         int            `json:"quantity" gorm:"not null"`
	Product          Product        `json:"product" gorm:"foreignKey:ProductID"`
	ProductVariant   ProductVariant `json:"product_variant" gorm:"foreignKey:ProductVariantID"`
}
