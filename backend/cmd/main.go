package main

import (
	"alexbirbirdev/go-shop/backend/config"
	"alexbirbirdev/go-shop/backend/internal/routes"

	"github.com/gin-gonic/gin"
)

func init() {
	config.GetEnv()
	config.InitDB()
}

func main() {
	r := gin.Default()
	routes.SetupRoutes(r)
	r.Run()
}
