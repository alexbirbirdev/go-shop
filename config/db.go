package config

import (
	"alexbirbirdev/go-shop/internal/models"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	var err error

	dsn := os.Getenv("DB_DSN")
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}

	DB.AutoMigrate(
		&models.Product{},
		&models.User{},
		&models.Category{},
		&models.ProductVariant{},
		&models.Favorite{},
		&models.CartItem{},
		&models.Order{},
		&models.OrderItem{},
		&models.UserAddress{},
	)

	return DB
}
