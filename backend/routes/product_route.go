package routes

import (
	"LearningCampusKabre/controllers"
	"LearningCampusKabre/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupProductRoutes(router *gin.Engine, db *gorm.DB) {
	productController := &controllers.ProductController{DB: db}

	productRoutes := router.Group("/api/products")
	{
		productRoutes.POST("", middlewares.AuthMiddleware(), middlewares.RequireRole("admin"), productController.CreateProduct)
		productRoutes.GET("", middlewares.AuthMiddleware(), productController.GetAllProducts)
		productRoutes.GET("/:id", middlewares.AuthMiddleware(), productController.GetProduct)
		productRoutes.PUT("/:id", middlewares.AuthMiddleware(), productController.UpdateProduct)
		productRoutes.DELETE("/:id", middlewares.AuthMiddleware(), middlewares.RequireRole("admin"), productController.DeleteProduct)
		productRoutes.DELETE("/softdelete/:id", middlewares.AuthMiddleware(), productController.SoftDeleteProduct)
	}
}
