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

// MenuInput représente les données attendues pour créer un menu
type MenuInput struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Items       []int   `json:"items"` // IDs des produits
	Price       float32 `json:"price"`
}

// ContainsID vérifie si un slice de produits contient un produit avec l'ID donné
func containsID(products []models.Product, id uint) bool {
	for _, p := range products {
		if p.ID == id {
			return true
		}
	}
	return false
}

// CreateMenu crée un nouveau menu
// @Summary Create a new menu
// @Description Create a new menu with associated products
// @Tags menus
// @Accept json
// @Produce json
// @Param menu body MenuInput true "Menu data"
// @Success 201 {object} models.Menu
// @Router /api/menus [post]
func (mc *MenuController) CreateMenu(c *gin.Context) {

	var request MenuInput

	// Vérification du JSON
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(request.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Le menu doit contenir au moins un produit"})
		return
	}

	// Vérifier que tous les produits existent et sont disponibles
	var products []models.Product

	if err := mc.DB.Where("id IN ? AND is_available = ?", request.Items, true).Find(&products).Error; err != nil {
		c.JSON(500, gin.H{"error": "Erreur lors de la vérification des produits"})
		return
	}

	if len(products) != len(request.Items) {
		c.JSON(400, gin.H{"error": "Un ou plusieurs produits ne sont pas disponibles"})
		return
	}

	// Création du menu
	menu := models.Menu{
		Name:        request.Name,
		Description: request.Description,
		Price:       request.Price,
	}

	if err := mc.DB.Create(&menu).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Création des MenuItems
	for _, p := range products {
		mc.DB.Create(&models.MenuItem{
			MenuID:      menu.ID,
			Name:        p.Name,
			Price:       p.Price,
			ImageURL:    p.ImageURL,
			Description: p.Description,
			Type:        p.Type,
		})
	}

	c.JSON(201, gin.H{
		"message": "Menu créé avec succès",
		"menu":    menu,
	})
}

// GetAllMenus récupère tous les menus
// @Summary Get all menus with their items
// @Description Retrieve a list of all menus with their associated items
// @Tags menus
// @Accept json
// @Produce json
// @Success 200 {array} models.Menu
// @Router /api/menus [get]
func (mc *MenuController) GetAllMenus(c *gin.Context) {
	var menus []models.Menu

	if err := mc.DB.
		Preload("MenuItems").
		Find(&menus).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des menus"})
		return
	}

	c.JSON(http.StatusOK, menus)
}
