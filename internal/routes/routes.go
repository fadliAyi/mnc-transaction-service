package routes

import (
	"transaction-service/internal/controllers"
	"transaction-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	router := gin.Default()

	router.POST("/register", controllers.Register)
	router.POST("/login", controllers.Login)

	// Authenticated routes
	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware())

	protected.POST("/topup", controllers.Topup)
	protected.POST("/pay", controllers.Payment)
	protected.POST("/transfer", controllers.Transfer)
	protected.GET("/transactions", controllers.TransactionReport)
	protected.PUT("/profile", controllers.UpdateProfile)

	return router
}
