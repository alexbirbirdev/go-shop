package models

import "gorm.io/gorm"

type ProductVariant struct {
	gorm.Model
	ProductID uint    `json:"product_id" gorm:"not null;index:idx_product_id_name,unique"`
	Name      string  `json:"name" gorm:"index:idx_product_id_name,unique"`
	Stock     int     `json:"stock" gorm:"not null"`
	Price     float64 `json:"price" gorm:"not null"`
}
