package models

import "gorm.io/gorm"

type UserAddress struct {
	gorm.Model
	UserID    uint   `json:"user_id" gorm:"not null"`
	Title     string `json:"title" gorm:"not null"`
	City      string `json:"city" gorm:"not null"`
	Street    string `json:"street" gorm:"not null"`
	House     string `json:"house" gorm:"not null"`
	Apartment string `json:"apartment"`
	Floor     string `json:"floor"`
	Comment   string `json:"comment"`
	IsDefault bool   `json:"is_default" gorm:"default:false"`
}
