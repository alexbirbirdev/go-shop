package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserID     uint        `json:"user_id" gorm:"not null"`
	Status     string      `json:"status" gorm:"not null"`
	TotalPrice float64     `json:"total_price" gorm:"not null"`
	OrderItems []OrderItem `json:"order_items" gorm:"foreignKey:OrderID"`
}
