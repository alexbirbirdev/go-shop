package models

import "gorm.io/gorm"

type OrderItem struct {
	gorm.Model
	OrderID          uint           `json:"order_id" gorm:"not null"`
	ProductID        uint           `json:"product_id" gorm:"not null"`
	ProductVariantID uint           `json:"product_variant_id" gorm:"not null"`
	Quantity         int            `json:"stock" gorm:"not null"`
	Price            float64        `json:"price" gorm:"not null"`
	Product          Product        `json:"product" gorm:"foreignKey:ProductID"`
	ProductVariant   ProductVariant `json:"product_variant" gorm:"foreignKey:ProductVariantID"`
}
