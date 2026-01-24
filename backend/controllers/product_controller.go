package controllers

//Appel de CreateProduct, GetAllProducts, GetProduct, UpdateProduct, DeleteProduct

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

// CreateProduct gère la création d'un produit
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

// GetAllProducts récupère tous les produits
func (pc *ProductController) GetAllProducts(c *gin.Context) {
	products, err := models.GetAllProducts(pc.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des produits"})
		return
	}

	c.JSON(http.StatusOK, products)
}

// GetProduct récupère un produit par son ID
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

// UpdateProduct met à jour un produit
func (pc *ProductController) UpdateProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	// Récupérer le produit existant
	product, err := models.GetProductByID(pc.DB, uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Produit non trouvé"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération du produit"})
		return
	}

	// Mettre à jour avec les nouvelles données
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Sauvegarder les modifications
	if err := models.UpdateProduct(pc.DB, product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour du produit"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// DeleteProduct supprime un produit
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
