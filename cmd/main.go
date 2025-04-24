package main

import (
	"alexbirbirdev/go-shop/config"
	"alexbirbirdev/go-shop/internal/routes"

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
