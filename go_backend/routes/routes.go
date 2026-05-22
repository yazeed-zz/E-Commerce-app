package routes

import (
	"ecommerce/backend/handlers"
	"ecommerce/backend/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")

	// Auth routes
	auth := api.Group("/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
	}

	// Protected routes
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/products", handlers.GetProducts)
		protected.GET("/products/:id", handlers.GetProduct)
		protected.POST("/products/:id/reviews", handlers.CreateReview)
		protected.GET("/products/:id/reviews", handlers.GetReviews)
		protected.POST("/cart", handlers.AddToCart)
		protected.GET("/cart", handlers.GetCart)
		protected.POST("/orders", handlers.CreateOrder)
		protected.GET("/orders", handlers.GetOrders)
		protected.GET("/recommendations/:type", handlers.GetRecommendations)
	}
	// Admin routes
	admin := api.Group("/admin")
	admin.Use(middleware.AuthMiddleware())
	{
		admin.GET("/stats", handlers.GetStats)
		admin.GET("/users", handlers.GetUsers)
		admin.POST("/products", handlers.CreateProduct)
		admin.DELETE("/products/:id", handlers.DeleteProduct)
		admin.PUT("/orders/:id/status", handlers.UpdateOrderStatus)
	}
}
