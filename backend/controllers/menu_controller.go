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

	// Récupérer les produits disponibles par type
	entrees, _ := models.GetAvailableItemsByType(mc.DB, "entree")
	plats, _ := models.GetAvailableItemsByType(mc.DB, "plat")
	desserts, _ := models.GetAvailableItemsByType(mc.DB, "dessert")

	// Structure attendue du front
	var request struct {
		EntreeID  uint `json:"entree_id"`
		PlatID    uint `json:"plat_id"`
		DessertID uint `json:"dessert_id"`
	}

	// Si aucun choix envoyé → renvoyer les options
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(200, gin.H{
			"message":  "Veuillez choisir une entrée, un plat et un dessert",
			"entrees":  entrees,
			"plats":    plats,
			"desserts": desserts,
		})
		return
	}

	// Vérifier que les IDs envoyés existent dans les produits disponibles
	if !containsID(entrees, request.EntreeID) ||
		!containsID(plats, request.PlatID) ||
		!containsID(desserts, request.DessertID) {

		println("Entree ID", request.EntreeID, "Plat ID", request.PlatID, "Dessert ID", request.DessertID)
		c.JSON(400, gin.H{
			"error": "Un des éléments sélectionnés n'est pas disponible",
		})
		return
	}

	// Construire le menu
	menu := models.Menu{
		EntreeID:  request.EntreeID,
		PlatID:    request.PlatID,
		DessertID: request.DessertID,
	}

	// Sauvegarder en base
	if err := mc.DB.Create(&menu).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Recharge le menu avec les relations
	mc.DB.
		Preload("Entree").
		Preload("Plat").
		Preload("Dessert").
		First(&menu, menu.ID)

	c.JSON(201, gin.H{
		"message": "Menu créé avec succès",
		"menu":    menu,
	})
}

func (mc *MenuController) GetAllMenus(c *gin.Context) {
	menus, err := models.GetAllMenus(mc.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des menus"})
		return
	}
	c.JSON(http.StatusOK, menus)
}
