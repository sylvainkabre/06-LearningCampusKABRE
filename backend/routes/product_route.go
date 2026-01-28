package routes

import (
	"LearningCampusKabre/controllers"
	"LearningCampusKabre/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupProductRoutes configure toutes les routes des produits
func SetupProductRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialiser le contrôleur
	productController := &controllers.ProductController{DB: db}

	// Groupe de routes pour les produits
	productRoutes := router.Group("/api/products")
	{
		productRoutes.POST("", middlewares.AuthMiddleware(), middlewares.AdminOnly(), productController.CreateProduct)
		productRoutes.GET("", middlewares.AuthMiddleware(), productController.GetAllProducts)
		productRoutes.GET("/:id", middlewares.AuthMiddleware(), productController.GetProduct)
		productRoutes.PUT("/:id", middlewares.AuthMiddleware(), productController.UpdateProduct)
		productRoutes.DELETE("/:id", middlewares.AuthMiddleware(), middlewares.AdminOnly(), productController.DeleteProduct)
		productRoutes.DELETE("/softdelete/:id", middlewares.AuthMiddleware(), productController.SoftDeleteProduct)
	}
}
