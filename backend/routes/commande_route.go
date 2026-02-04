package routes

import (
	"LearningCampusKabre/controllers"
	"LearningCampusKabre/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupCommandesRoutes configure toutes les routes des commandes
func SetupCommandesRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialiser le contrôleur
	commandeController := &controllers.CommandeController{DB: db}

	// Groupe de routes pour les commandes
	commandeRoutes := router.Group("/api/commandes")
	{
		commandeRoutes.POST("", middlewares.AuthMiddleware(), commandeController.CreateCommande)
		commandeRoutes.GET("", middlewares.AuthMiddleware(), commandeController.GetAllCommandes)
		commandeRoutes.GET("/:id", middlewares.AuthMiddleware(), commandeController.GetCommandeByID)
		commandeRoutes.PUT("/:id", middlewares.AuthMiddleware(), commandeController.UpdateCommande)
		commandeRoutes.DELETE("/:id", middlewares.AuthMiddleware(), commandeController.DeleteCommande)
	}
}
