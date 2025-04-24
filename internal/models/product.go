package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name        string           `json:"name" gorm:"not null"`
	Description string           `json:"description" gorm:"not null"`
	Price       float64          `json:"price" gorm:"not null"`
	Stock       int              `json:"stock" gorm:"not null"`
	CategoryID  uint             `json:"category_id" gorm:"not null"`
	Image       string           `json:"image" gorm:"not null"`
	Variants    []ProductVariant `json:"variants" gorm:"foreignKey:ProductID"`
}
