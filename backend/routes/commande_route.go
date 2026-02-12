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
		commandeRoutes.POST("", middlewares.AuthMiddleware(), middlewares.RequireRole("admin", "receiver", "preparer"), commandeController.CreateCommande)
		commandeRoutes.GET("", middlewares.AuthMiddleware(), commandeController.GetAllCommandes)
		commandeRoutes.GET("/:id", middlewares.AuthMiddleware(), commandeController.GetCommandeByID)
		commandeRoutes.PUT("/admin/:id", middlewares.AuthMiddleware(), middlewares.RequireRole("admin"), commandeController.AdminUpdateCommande)
		commandeRoutes.PUT("/preparer/:id", middlewares.AuthMiddleware(), middlewares.RequireRole("preparer"), commandeController.PreparerUpdateCommande)
		commandeRoutes.PUT("/receiver/:id", middlewares.AuthMiddleware(), middlewares.RequireRole("receiver"), commandeController.ReceiverUpdateCommande)
		commandeRoutes.DELETE("/:id", middlewares.AuthMiddleware(), middlewares.RequireRole("admin", "receiver", "preparer"), commandeController.DeleteCommande)
	}
}

// CRUD regarder définition
// Ajouter routes Update propre à chaque utilisateur
