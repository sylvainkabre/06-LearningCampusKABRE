package routes

import (
	"LearningCampusKabre/controllers"
	"LearningCampusKabre/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupMenuRoutes(router *gin.Engine, db *gorm.DB) {
	menuController := &controllers.MenuController{DB: db}

	menuRoutes := router.Group("/api/menus")
	{
		menuRoutes.POST("", middlewares.AuthMiddleware(), middlewares.RequireRole("admin"), menuController.CreateMenu)
		menuRoutes.GET("", menuController.GetAllMenus)
		menuRoutes.GET("/:id", middlewares.AuthMiddleware(), menuController.GetMenuByID)
		menuRoutes.PUT("/:id", middlewares.AuthMiddleware(), menuController.UpdateMenu)
		menuRoutes.DELETE("/:id", middlewares.AuthMiddleware(), middlewares.RequireRole("admin"), menuController.DeleteMenu)
		menuRoutes.DELETE("/softdelete/:id", middlewares.AuthMiddleware(), menuController.SoftDeleteMenu)
	}
}
