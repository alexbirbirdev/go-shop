package main

import (
	"alexbirbirdev/go-shop/backend/config"
	"alexbirbirdev/go-shop/backend/internal/middleware"
	"alexbirbirdev/go-shop/backend/internal/routes"

	"github.com/gin-gonic/gin"
)

func init() {
	config.GetEnv()
	config.InitDB()
}

func main() {
	r := gin.Default()

	// r.Use(cors.New(cors.Config{
	// 	AllowOrigins: []string{"http://localhost:5173"},
	// 	AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	// 	AllowHeaders: []string{
	// 		"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With",
	// 	},
	// 	AllowCredentials: true,
	// }))
	// r.OPTIONS("/*path", func(c *gin.Context) {
	// 	c.Status(http.StatusOK)
	// })

	r.Use(middleware.CORS())
	r.Static("/uploads", "./uploads")

	routes.SetupRoutes(r)
	r.Run()
}
