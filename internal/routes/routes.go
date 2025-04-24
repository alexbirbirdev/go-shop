package routes

import (
	"alexbirbirdev/go-shop/internal/handlers"
	"alexbirbirdev/go-shop/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	productRoutes := r.Group("/products")
	{
		productRoutes.POST("/", handlers.CreateProduct)
		productRoutes.GET("/", handlers.GetProducts)
		productRoutes.GET("/:id", handlers.GetProduct)
		productRoutes.DELETE("/:id", handlers.DeleteProduct)
		productRoutes.PUT("/:id", handlers.UpdateProduct)

		productRoutes.POST("/:id/variants", handlers.CreateProductVariant)
		productRoutes.GET("/:id/variants", handlers.GetProductVariants)
	}

	categoryRoutes := r.Group("/category")
	{
		categoryRoutes.POST("/", handlers.CreateCategory)
		categoryRoutes.GET("/", handlers.GetCategories)
		categoryRoutes.GET("/:id", handlers.GetCategory)
		categoryRoutes.DELETE("/:id", handlers.DeleteCategory)
		categoryRoutes.PUT("/:id", handlers.UpdateCategory)
	}

	productVariantRoutes := r.Group("/variants")
	{
		productVariantRoutes.GET("/:id", handlers.GetProductVariant)
		productVariantRoutes.DELETE("/:id", handlers.DeleteProductVariant)
		productVariantRoutes.PUT("/:id", handlers.UpdateProductVariant)
	}

	cartItemRoutes := r.Group("/cart")
	{
		cartItemRoutes.POST("/", handlers.CreateCartItem)
		cartItemRoutes.GET("/", handlers.GetCartItems)
		cartItemRoutes.DELETE("/:id", handlers.DeleteCartItem)
		cartItemRoutes.DELETE("/clear", handlers.ClearCart)
		cartItemRoutes.PUT("/:id", handlers.UpdateQuantity)
		cartItemRoutes.POST("/:id/increment", handlers.IncrementCartItem)
		cartItemRoutes.POST("/:id/decrement", handlers.DecrementCartItem)
	}

	orderRoutes := r.Group("/orders")
	{
		orderRoutes.POST("/", handlers.CreateOrder)
		orderRoutes.GET("/", handlers.GetOrders)
		orderRoutes.GET("/:id", handlers.GetOrder)
	}

	userAddressesRoutes := r.Group("/user-addresses").Use(middleware.AuthMiddleware())
	{
		userAddressesRoutes.POST("/", handlers.CreateUserAddresses)
		userAddressesRoutes.GET("/", handlers.GetUserAddresses)
		userAddressesRoutes.GET("/:id", handlers.GetUserAddress)
		userAddressesRoutes.DELETE("/:id", handlers.DeleteUserAddress)
		userAddressesRoutes.PUT("/:id", handlers.UpdateUserAddress)
	}

	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/signin", middleware.AuthRequired, handlers.SignIn)
		authRoutes.POST("/signup", middleware.AuthRequired, handlers.SignUp)
		authRoutes.POST("/logout", middleware.AuthMiddleware(), handlers.Logout)
	}
}
