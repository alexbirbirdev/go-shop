package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name     string     `json:"name" gorm:"not null"`
	ParentID *uint      `json:"parent_id" gorm:"default:NULL"`
	Children []Category `json:"children" gorm:"foreignKey:ParentID"`
}
