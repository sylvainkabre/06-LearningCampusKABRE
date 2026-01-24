package routes

import (
	"LearningCampusKabre/controllers"

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
		productRoutes.POST("", productController.CreateProduct)
		productRoutes.GET("", productController.GetAllProducts)
		productRoutes.GET("/:id", productController.GetProduct)
		productRoutes.PUT("/:id", productController.UpdateProduct)
		productRoutes.DELETE("/:id", productController.DeleteProduct)
	}
}
