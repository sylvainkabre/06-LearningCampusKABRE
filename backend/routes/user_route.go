package routes

import (
	"LearningCampusKabre/controllers"
	"LearningCampusKabre/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupUserRoutes(router *gin.Engine, db *gorm.DB) {
	userController := controllers.RefUserController(db)
	authController := controllers.RefAuthController(db)

	auth := router.Group("/api/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
	}

	users := router.Group("/api/users")
	users.Use(middlewares.AuthMiddleware())
	{
		users.POST("", userController.CreateUser)
		users.GET("", userController.GetAllUsers)
		users.GET("/:id", userController.GetUserByID)
		users.PUT("/:id", userController.UpdateUser)
		users.DELETE("/:id", middlewares.RequireRole("admin"), userController.DeleteUser)
	}
}
