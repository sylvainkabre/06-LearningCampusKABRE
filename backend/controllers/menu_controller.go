package controllers

import (
	"LearningCampusKabre/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type MenuController struct {
	DB *gorm.DB
}

func RefMenuController(db *gorm.DB) *MenuController {
	return &MenuController{DB: db}
}

// MenuInput représente les données attendues pour créer un menu
type MenuInput struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Items       []int           `json:"items"`
	Price       decimal.Decimal `json:"price"`
	ImageURL    string          `json:"image_url"`
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
// @Router /menus [post]
// @Security BearerAuth
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
		ImageURL:    request.ImageURL,
	}

	if err := mc.DB.Create(&menu).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Création des MenuItems
	for _, p := range products {
		println(p.Price.String())
		mc.DB.Create(&models.MenuItem{
			MenuID:      menu.ID,
			Name:        p.Name,
			Price:       p.Price,
			ImageURL:    p.ImageURL,
			Description: p.Description,
			Type:        p.Type,
		})
	}

	// Il me semblait que Preload pouvait être chaîné après une création, mais apparemment non
	mc.DB.Preload("MenuItems").First(&menu, menu.ID)

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
// @Router /menus [get]
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

// GetMenuByID récupère un menu par son ID
// @Summary Get a menu by ID
// @Description Retrieve a single menu with its items by ID
// @Tags menus
// @Produce json
// @Param id path int true "Menu ID"
// @Success 200 {object} models.Menu
// @Failure 404 {object} map[string]string
// @Router /menus/{id} [get]
// @Security BearerAuth
func (mc *MenuController) GetMenuByID(c *gin.Context) {
	id := c.Param("id")

	var menu models.Menu
	if err := mc.DB.Preload("MenuItems").First(&menu, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Menu non trouvé"})
		return
	}

	c.JSON(http.StatusOK, menu)
}

// UpdateMenu met à jour un menu existant
// @Summary Update a menu
// @Description Update an existing menu and its associated products
// @Tags menus
// @Accept json
// @Produce json
// @Param id path int true "Menu ID"
// @Param menu body MenuInput true "Menu data"
// @Success 200 {object} models.Menu
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /menus/{id} [put]
// @Security BearerAuth
func (mc *MenuController) UpdateMenu(c *gin.Context) {
	id := c.Param("id")

	var menu models.Menu
	if err := mc.DB.First(&menu, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Menu non trouvé"})
		return
	}

	var request MenuInput
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(request.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Le menu doit contenir au moins un produit"})
		return
	}

	var products []models.Product
	if err := mc.DB.Where("id IN ? AND is_available = ?", request.Items, true).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la vérification des produits"})
		return
	}

	if len(products) != len(request.Items) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Un ou plusieurs produits ne sont pas disponibles"})
		return
	}

	// Mise à jour des champs
	menu.Name = request.Name
	menu.Description = request.Description
	menu.Price = request.Price
	menu.ImageURL = request.ImageURL

	if err := mc.DB.Save(&menu).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour du menu"})
		return
	}

	// Suppression des anciens MenuItems et recréation
	mc.DB.Where("menu_id = ?", menu.ID).Delete(&models.MenuItem{})
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

	mc.DB.Preload("MenuItems").First(&menu, menu.ID)
	c.JSON(http.StatusOK, menu)
}

// DeleteMenu supprime définitivement un menu (hard delete)
// @Summary Hard delete a menu
// @Description Permanently delete a menu and its items
// @Tags menus
// @Param id path int true "Menu ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /menus/{id} [delete]
// @Security BearerAuth
func (mc *MenuController) DeleteMenu(c *gin.Context) {
	id := c.Param("id")

	var menu models.Menu
	if err := mc.DB.First(&menu, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Menu non trouvé"})
		return
	}

	mc.DB.Where("menu_id = ?", menu.ID).Delete(&models.MenuItem{})

	if err := mc.DB.Unscoped().Delete(&menu).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression du menu"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Menu supprimé définitivement"})
}

// SoftDeleteMenu effectue un soft delete d'un menu
// @Summary Soft delete a menu
// @Description Soft delete a menu (marks as deleted without removing from DB)
// @Tags menus
// @Param id path int true "Menu ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /menus/softdelete/{id} [delete]
// @Security BearerAuth
func (mc *MenuController) SoftDeleteMenu(c *gin.Context) {
	id := c.Param("id")

	var menu models.Menu
	if err := mc.DB.First(&menu, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Menu non trouvé"})
		return
	}

	if err := mc.DB.Delete(&menu).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du soft delete du menu"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Menu soft supprimé"})
}
