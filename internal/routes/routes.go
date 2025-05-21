package routes

import (
	"alexbirbirdev/go-shop/internal/handlers"
	"alexbirbirdev/go-shop/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	productRoutes := r.Group("/products")
	{
		productRoutes.GET("/", handlers.GetProducts)
		productRoutes.GET("/:id", handlers.GetProduct)
		productRoutes.GET("/:id/variants", handlers.GetProductVariants)
	}

	categoryRoutes := r.Group("/category")
	{
		categoryRoutes.GET("/", handlers.GetCategories)
		categoryRoutes.GET("/:id", handlers.GetCategory)
	}

	productVariantRoutes := r.Group("/variants")
	{
		productVariantRoutes.GET("/:id", handlers.GetProductVariant)
	}

	cartItemRoutes := r.Group("/cart").Use(middleware.AuthMiddleware())
	{
		cartItemRoutes.POST("/", handlers.CreateCartItem)
		cartItemRoutes.GET("/", handlers.GetCartItems)
		cartItemRoutes.DELETE("/:id", handlers.DeleteCartItem)
		cartItemRoutes.DELETE("/clear", handlers.ClearCart)
		cartItemRoutes.PUT("/:id", handlers.UpdateQuantity)
		cartItemRoutes.POST("/:id/increment", handlers.IncrementCartItem)
		cartItemRoutes.POST("/:id/decrement", handlers.DecrementCartItem)
		cartItemRoutes.GET("/:id/check", handlers.CheckInCart)
	}

	orderRoutes := r.Group("/orders").Use(middleware.AuthMiddleware())
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

	userProfileRoutes := r.Group("/profile").Use(middleware.AuthMiddleware())
	{
		userProfileRoutes.GET("/", handlers.GetUserProfile)
		userProfileRoutes.PUT("/", handlers.UpdateUserProfile)
	}

	favoriteRoutes := r.Group("/favorites").Use(middleware.AuthMiddleware())
	{
		favoriteRoutes.POST("/", handlers.AddFavorite)
		favoriteRoutes.GET("/", handlers.GetFavorites)
		favoriteRoutes.GET("/:id", handlers.CheckFavorite)
		favoriteRoutes.DELETE("/:id", handlers.DeleteFavorite)
		favoriteRoutes.DELETE("/", handlers.ClearFavorites)
	}

	// Возможности админа
	adminRoutes := r.Group("/admin")
	adminRoutes.Use(middleware.AuthMiddleware(), middleware.AdminOnly())
	adminProductRoutes := adminRoutes.Group("/products")
	{
		adminProductRoutes.GET("/", handlers.AdminGetProducts)
		adminProductRoutes.GET("/:id", handlers.AdminGetProduct)
		adminProductRoutes.DELETE("/:id", handlers.DeleteProduct)
		adminProductRoutes.POST("/", handlers.CreateProduct)
		adminProductRoutes.PUT("/:id", handlers.UpdateProduct)
		adminProductRoutes.POST("/:id/variants", handlers.CreateProductVariant)
	}
	adminProductVariantRoutes := adminRoutes.Group("/variants")
	{
		adminProductVariantRoutes.DELETE("/:id", handlers.DeleteProductVariant)
		adminProductVariantRoutes.PUT("/:id", handlers.UpdateProductVariant)
	}
	adminCategoryRoutes := adminRoutes.Group("/category")
	{
		adminCategoryRoutes.DELETE("/:id", handlers.DeleteCategory)
		adminCategoryRoutes.PUT("/:id", handlers.UpdateCategory)
		adminCategoryRoutes.POST("/", handlers.CreateCategory)
		adminCategoryRoutes.GET("/", handlers.AdminGetCategories)
		adminCategoryRoutes.GET("/:id", handlers.AdminGetCategory)
		adminCategoryRoutes.GET("/all", handlers.ShowAllCategories)
	}
	adminOrdersRoutes := adminRoutes.Group("/orders")
	{
		adminOrdersRoutes.GET("/", handlers.AdminGetOrders)
		adminOrdersRoutes.GET("/:id", handlers.AdminGetOrder)
		adminOrdersRoutes.PUT("/:id/status", handlers.UpdateOrderStatus)
	}
	adminUsersRoutes := adminRoutes.Group("/users")
	{
		adminUsersRoutes.GET("/", handlers.AdminGetUsers)
		adminUsersRoutes.GET("/:id", handlers.AdminGetUser)
		adminUsersRoutes.PUT("/:id", handlers.AdminUpdateUser)
	}
}
