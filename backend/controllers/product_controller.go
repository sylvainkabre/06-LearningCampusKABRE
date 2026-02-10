package controllers

import (
	"net/http"
	"strconv"

	"LearningCampusKabre/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductController struct {
	DB *gorm.DB
}

func RefProductController(db *gorm.DB) *ProductController {
	return &ProductController{DB: db}
}

// CreateProduct
// @Summary Create a new product
// @Description Create a new product and return it
// @Tags products
// @Accept json
// @Produce json
// @Param product body models.Product true "Product data"
// @Success 201 {object} models.Product
// @Router /api/products [post]
func (pc *ProductController) CreateProduct(c *gin.Context) {
	var product models.Product

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.CreateProduct(pc.DB, &product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création du produit"})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// GetAllProducts
// @Summary Get all products
// @Description Retrieve a list of all products
// @Tags products
// @Produce json
// @Success 200 {array} models.Product
// @Router /products [get]
func (pc *ProductController) GetAllProducts(c *gin.Context) {
	products, err := models.GetAllProducts(pc.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des produits"})
		return
	}

	c.JSON(http.StatusOK, products)
}

// GetProduct
// @Summary Get a product by ID
// @Description Retrieve a single product by its ID
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} models.Product
// @Failure 404 {object} map[string]string
// @Router /products/{id} [get]
func (pc *ProductController) GetProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	product, err := models.GetProductByID(pc.DB, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// UpdateProduct
// @Summary Update a product
// @Description Update an existing product by its ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body models.Product true "Updated product data"
// @Success 200 {object} models.Product
// @Router /products/{id} [put]
func (pc *ProductController) UpdateProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	product, err := models.GetProductByID(pc.DB, uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Produit non trouvé"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération du produit"})
		return
	}

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.UpdateProduct(pc.DB, product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour du produit"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// DeleteProduct
// @Summary Delete a product
// @Description Delete an existing product by its ID
// @Tags products
// @Param id path int true "Product ID"
// @Success 200 {object} map[string]string
// @Router /products/{id} [delete]
func (pc *ProductController) DeleteProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	if err := models.DeleteProduct(pc.DB, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression du produit"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Produit supprimé avec succès"})
}

// SoftDeleteProduct
// @Summary Soft delete a product
// @Description Soft delete an existing product by its ID (sets deleted_at)
// @Tags products
// @Param id path int true "Product ID"
// @Success 200 {object} map[string]string
// @Router /products/{id}/soft [delete]
func (ctrl *ProductController) SoftDeleteProduct(c *gin.Context) {
	id := c.Param("id")

	var product models.Product

	if err := ctrl.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produit non trouvé"})
		return
	}

	if err := ctrl.DB.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du soft delete"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Produit soft supprimé"})
}
