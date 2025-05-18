package models

import "gorm.io/gorm"

type Favorite struct {
	gorm.Model
	UserID    uint    `json:"user_id" gorm:"not null"`
	ProductID uint    `json:"product_id" gorm:"not null"`
	Product   Product `json:"product" gorm:"foreignKey:ProductID;references:ID"`
}
