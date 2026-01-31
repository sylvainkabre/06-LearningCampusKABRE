package controllers

import (
	"LearningCampusKabre/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MenuController struct {
	DB *gorm.DB
}

// Vérifie si un ID appartient à une liste de produits disponibles
func containsID(products []models.Product, id uint) bool {
	for _, p := range products {
		if p.ID == id {
			return true
		}
	}
	return false
}

func (mc *MenuController) CreateMenu(c *gin.Context) {

	// Structure attendue du front
	var request struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Items       []struct {
			ProductID uint `json:"product_id"`
		} `json:"items"`
	}

	// Vérification du JSON
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format JSON invalide"})
		return
	}

	if len(request.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Le menu doit contenir au moins un produit"})
		return
	}

	// Vérifier que tous les produits existent et sont disponibles
	var products []models.Product
	var productIDs []uint

	for _, item := range request.Items {
		productIDs = append(productIDs, item.ProductID)
	}

	if err := mc.DB.Where("id IN ? AND is_available = ?", productIDs, true).Find(&products).Error; err != nil {
		c.JSON(500, gin.H{"error": "Erreur lors de la vérification des produits"})
		return
	}

	if len(products) != len(request.Items) {
		c.JSON(400, gin.H{"error": "Un ou plusieurs produits ne sont pas disponibles"}) // On pourrait être plus précis ici
		return
	}

	// Calcul du prix total
	var total float32
	for _, p := range products {
		total += p.Price
	}

	// Création du menu
	menu := models.Menu{
		Name:        request.Name,
		Description: request.Description,
		Price:       total,
	}

	if err := mc.DB.Create(&menu).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()}) //Idem on pourrait être plus précis
		return
	}

	// Création des MenuItems
	for _, p := range products {
		mc.DB.Create(&models.MenuItem{
			MenuID:    menu.ID,
			ProductID: p.ID,
		})
	}

	// Recharger le menu complet
	mc.DB.
		Preload("MenuItems.Product").
		First(&menu, menu.ID)

	c.JSON(201, gin.H{
		"message": "Menu créé avec succès",
		"menu":    menu,
	})
}

func (mc *MenuController) GetAllMenus(c *gin.Context) {
	var menus []models.Menu

	if err := mc.DB.
		Preload("MenuItems.Product").
		Find(&menus).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des menus"})
		return
	}

	c.JSON(http.StatusOK, menus)
}
