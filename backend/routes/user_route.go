package routes

import (
	"LearningCampusKabre/controllers"
	"LearningCampusKabre/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupUserRoutes(router *gin.Engine, db *gorm.DB) {
	userController := controllers.NewUserController(db)

	userGroup := router.Group("/api/users")

	{
		userGroup.POST("", middlewares.AuthMiddleware(), userController.CreateUser)
		userGroup.GET("", middlewares.AuthMiddleware(), userController.GetAllUsers)
		userGroup.GET("/:id", middlewares.AuthMiddleware(), userController.GetUserByID)
		userGroup.PUT("/:id", middlewares.AuthMiddleware(), userController.UpdateUser)
		userGroup.DELETE("/:id", middlewares.AuthMiddleware(), userController.DeleteUser)
		userGroup.POST("/register", controllers.Register)
		userGroup.POST("/login", controllers.Login)

	}
}
