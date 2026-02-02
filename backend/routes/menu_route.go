package routes

import (
	"LearningCampusKabre/controllers"
	"LearningCampusKabre/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupMenuRoutes configure toutes les routes des menus
func SetupMenuRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialiser le contrôleur
	menuController := &controllers.MenuController{DB: db}

	// Groupe de routes pour les menus
	menuRoutes := router.Group("/api/menus")
	{
		menuRoutes.POST("", middlewares.AuthMiddleware(), middlewares.AdminOnly(), menuController.CreateMenu)
		menuRoutes.GET("", menuController.GetAllMenus)
		//menuRoutes.GET("/:id", middlewares.AuthMiddleware(), menuController.GetMenu)
		//menuRoutes.PUT("/:id", middlewares.AuthMiddleware(), menuController.UpdateMenu)
		//menuRoutes.DELETE("/:id", middlewares.AuthMiddleware(), middlewares.AdminOnly(), menuController.DeleteMenu)
		//menuRoutes.DELETE("/softdelete/:id", middlewares.AuthMiddleware(), menuController.SoftDeleteMenu)
	}
}
