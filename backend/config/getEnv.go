package config

import (
	"fmt"
	"os"
)

func GetEnv() {
	fmt.Printf("DB_HOST: %s\n", os.Getenv("DB_HOST"))
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Println("Error loading .env file")
	// }
}
