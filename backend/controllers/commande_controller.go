package controllers

import (
	"LearningCampusKabre/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type CommandeController struct {
	DB *gorm.DB
}

type CommandeInput struct {
	Menus    []int           `json:"menus"`
	Products []int           `json:"products"`
	Price    decimal.Decimal `json:"price"`
}

// Struct utilisée pour la mise à jour d'une commande existante
type CommandeUpdateInput struct {
	Menus    []int             `json:"menus"`
	Products []int             `json:"products"`
	Price    decimal.Decimal   `json:"price"`
	Status   models.StatusType `json:"status"`
}

// CreateCommande crée une nouvelle commande
// @Summary Create a new commande
// @Description Create a new commande with associated menus and products
// @Tags commandes
// @Accept json
// @Produce json
// @Success 201 {object} models.Commande
// @Router /api/commandes [post]
func (cc *CommandeController) CreateCommande(c *gin.Context) {

	var request CommandeInput

	// Vérification du JSON
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(request.Menus) == 0 && len(request.Products) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "La commande doit contenir au moins un menu ou un produit"})
		return
	}

	// Vérifier que les menus existent
	var menus []models.Menu
	if len(request.Menus) > 0 {
		if err := cc.DB.Where("id IN ?", request.Menus).Find(&menus).Error; err != nil {
			c.JSON(500, gin.H{"error": "Erreur lors de la vérification des menus"})
			return
		}
		if len(menus) != len(request.Menus) {
			c.JSON(400, gin.H{"error": "Un ou plusieurs menus sont introuvables"})
			return
		}
	}

	// Vérifier que les produits existent
	var products []models.Product
	if len(request.Products) > 0 {
		if err := cc.DB.Where("id IN ?", request.Products).Find(&products).Error; err != nil {
			c.JSON(500, gin.H{"error": "Erreur lors de la vérification des produits"})
			return
		}
		if len(products) != len(request.Products) {
			c.JSON(400, gin.H{"error": "Un ou plusieurs produits sont introuvables"})
			return
		}
	}

	// Création de la commande et on la passe à l'état "pending"
	commande := models.Commande{
		Status: models.StatusPending,
		Price:  request.Price,
	}

	if err := cc.DB.Create(&commande).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Création des CommandeMenu (snapshots)
	for _, m := range menus {
		cc.DB.Create(&models.CommandeMenu{
			CommandeID:  commande.ID,
			MenuID:      m.ID,
			Name:        m.Name,
			Price:       m.Price,
			Description: m.Description,
			ImageURL:    m.ImageURL,
		})
	}

	// Création des CommandeProduct (snapshots)
	for _, p := range products {
		cc.DB.Create(&models.CommandeProduct{
			CommandeID:  commande.ID,
			ProductID:   p.ID,
			Name:        p.Name,
			Price:       p.Price,
			ImageURL:    p.ImageURL,
			Description: p.Description,
			Type:        p.Type,
		})
	}

	// Recharger la commande complète
	cc.DB.
		Preload("Menus").
		Preload("Products").
		First(&commande)

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Commande créée avec succès",
		"commande": commande,
	})
}

// GetAllCommandes récupère toutes les commandes
// @Summary Get all commandes
// @Description Get all commandes with their associated menus and products
// @Tags commandes
// @Accept json
// @Produce json
// @Success 201 {object} models.Commande
// @Router /api/commandes [get]

func (cc *CommandeController) GetAllCommandes(c *gin.Context) {

	commandes, err :=
		models.GetAllComm(cc.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des commandes"})
		return
	}
	c.JSON(http.StatusOK, commandes)
}

// GetCommandeByID récupère une commande par son ID
// @Summary Get a commande by ID
// @Description Get a commande by its ID with associated menus and products
// @Tags commandes
// @Accept json
// @Produce json
// @Success 201 {object} models.Commande
// @Router /api/commandes/{id} [get]
func (cc *CommandeController) GetCommandeByID(c *gin.Context) {
	idParam := c.Param("id")
	var id uint
	if _, err := fmt.Sscan(idParam, &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}
	commande, err := models.GetCommandeById(cc.DB, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, commande)
}

// UpdateCommande met à jour une commande existante
// @Summary Update an existing commande
// @Description Update an existing commande with new data
// @Tags commandes
// @Accept json
// @Produce json
// @Success 201 {object} models.Commande
// @Router /api/commandes/{id} [put]
func (cc *CommandeController) AdminUpdateCommande(c *gin.Context) {

	id := c.Param("id")

	var commande models.Commande
	if err := cc.DB.First(&commande, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Commande introuvable"})
		return
	}

	var request CommandeUpdateInput
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var menus []models.Menu
	if len(request.Menus) > 0 {
		if err := cc.DB.Where("id IN ?", request.Menus).Find(&menus).Error; err != nil {
			c.JSON(500, gin.H{"error": "Erreur lors de la vérification des menus"})
			return
		}
		if len(menus) != len(request.Menus) {
			c.JSON(400, gin.H{"error": "Un ou plusieurs menus sont introuvables"})
			return
		}
	}

	var products []models.Product
	if len(request.Products) > 0 {
		if err := cc.DB.Where("id IN ?", request.Products).Find(&products).Error; err != nil {
			c.JSON(500, gin.H{"error": "Erreur lors de la vérification des produits"})
			return
		}
		if len(products) != len(request.Products) {
			c.JSON(400, gin.H{"error": "Un ou plusieurs produits sont introuvables"})
			return
		}
	}

	commande.Price = request.Price
	commande.Status = request.Status
	cc.DB.Save(&commande)

	cc.DB.Where("commande_id = ?", commande.ID).Delete(&models.CommandeMenu{})
	cc.DB.Where("commande_id = ?", commande.ID).Delete(&models.CommandeProduct{})

	for _, m := range menus {
		cc.DB.Create(&models.CommandeMenu{
			CommandeID:  commande.ID,
			MenuID:      m.ID,
			Name:        m.Name,
			Price:       m.Price,
			Description: m.Description,
			ImageURL:    m.ImageURL,
		})
	}

	for _, p := range products {
		cc.DB.Create(&models.CommandeProduct{
			CommandeID:  commande.ID,
			ProductID:   p.ID,
			Name:        p.Name,
			Price:       p.Price,
			ImageURL:    p.ImageURL,
			Description: p.Description,
			Type:        p.Type,
		})
	}

	cc.DB.
		Preload("Menus").
		Preload("Products").
		First(&commande)

	c.JSON(http.StatusOK, gin.H{
		"message":  "Commande mise à jour avec succès",
		"commande": commande,
	})
}

// PreparerUpdateCommande met à jour une commande existante
// @Summary Update an existing commande to ready status
// @Description Update an existing commande with new data
// @Tags commandes
// @Accept json
// @Produce json
// @Success 201 {object} models.Commande
// @Router /api/commandes/{id} [put]
func (cc *CommandeController) PreparerUpdateCommande(c *gin.Context) {

	id := c.Param("id")

	var commande models.Commande
	if err := cc.DB.First(&commande, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Commande introuvable"})
		return
	}

	var request CommandeUpdateInput
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// On s'assure que le preparer passe le status à Ready et pas autre chose.
	if request.Status != "ready" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status non autorisé pour l'utilisateur preparer"})
		return
	}

	var menus []models.Menu
	if len(request.Menus) > 0 {
		if err := cc.DB.Where("id IN ?", request.Menus).Find(&menus).Error; err != nil {
			c.JSON(500, gin.H{"error": "Erreur lors de la vérification des menus"})
			return
		}
		if len(menus) != len(request.Menus) {
			c.JSON(400, gin.H{"error": "Un ou plusieurs menus sont introuvables"})
			return
		}
	}

	var products []models.Product
	if len(request.Products) > 0 {
		if err := cc.DB.Where("id IN ?", request.Products).Find(&products).Error; err != nil {
			c.JSON(500, gin.H{"error": "Erreur lors de la vérification des produits"})
			return
		}
		if len(products) != len(request.Products) {
			c.JSON(400, gin.H{"error": "Un ou plusieurs produits sont introuvables"})
			return
		}
	}

	commande.Price = request.Price
	commande.Status = request.Status
	cc.DB.Save(&commande)

	cc.DB.Where("commande_id = ?", commande.ID).Delete(&models.CommandeMenu{})
	cc.DB.Where("commande_id = ?", commande.ID).Delete(&models.CommandeProduct{})

	for _, m := range menus {
		cc.DB.Create(&models.CommandeMenu{
			CommandeID:  commande.ID,
			MenuID:      m.ID,
			Name:        m.Name,
			Price:       m.Price,
			Description: m.Description,
			ImageURL:    m.ImageURL,
		})
	}

	for _, p := range products {
		cc.DB.Create(&models.CommandeProduct{
			CommandeID:  commande.ID,
			ProductID:   p.ID,
			Name:        p.Name,
			Price:       p.Price,
			ImageURL:    p.ImageURL,
			Description: p.Description,
			Type:        p.Type,
		})
	}

	cc.DB.
		Preload("Menus").
		Preload("Products").
		First(&commande)

	c.JSON(http.StatusOK, gin.H{
		"message":  "Commande mise à jour avec succès",
		"commande": commande,
	})
}

// ReceiverUpdateCommande met à jour une commande existante to pending or delivered
// @Summary Update an existing commande to ready status
// @Description Update an existing commande with new data
// @Tags commandes
// @Accept json
// @Produce json
// @Success 201 {object} models.Commande
// @Router /api/commandes/{id} [put]
func (cc *CommandeController) ReceiverUpdateCommande(c *gin.Context) {

	id := c.Param("id")

	var commande models.Commande
	if err := cc.DB.First(&commande, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Commande introuvable"})
		return
	}

	var request CommandeUpdateInput
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// On s'assure que le preparer passe le status à Ready et pas autre chose.
	if request.Status != "pending" && request.Status != "delivered" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status non autorisé pour l'utilisateur receiver"})
		return
	}

	var menus []models.Menu
	if len(request.Menus) > 0 {
		if err := cc.DB.Where("id IN ?", request.Menus).Find(&menus).Error; err != nil {
			c.JSON(500, gin.H{"error": "Erreur lors de la vérification des menus"})
			return
		}
		if len(menus) != len(request.Menus) {
			c.JSON(400, gin.H{"error": "Un ou plusieurs menus sont introuvables"})
			return
		}
	}

	var products []models.Product
	if len(request.Products) > 0 {
		if err := cc.DB.Where("id IN ?", request.Products).Find(&products).Error; err != nil {
			c.JSON(500, gin.H{"error": "Erreur lors de la vérification des produits"})
			return
		}
		if len(products) != len(request.Products) {
			c.JSON(400, gin.H{"error": "Un ou plusieurs produits sont introuvables"})
			return
		}
	}

	commande.Price = request.Price
	commande.Status = request.Status
	cc.DB.Save(&commande)

	cc.DB.Where("commande_id = ?", commande.ID).Delete(&models.CommandeMenu{})
	cc.DB.Where("commande_id = ?", commande.ID).Delete(&models.CommandeProduct{})

	for _, m := range menus {
		cc.DB.Create(&models.CommandeMenu{
			CommandeID:  commande.ID,
			MenuID:      m.ID,
			Name:        m.Name,
			Price:       m.Price,
			Description: m.Description,
			ImageURL:    m.ImageURL,
		})
	}

	for _, p := range products {
		cc.DB.Create(&models.CommandeProduct{
			CommandeID:  commande.ID,
			ProductID:   p.ID,
			Name:        p.Name,
			Price:       p.Price,
			ImageURL:    p.ImageURL,
			Description: p.Description,
			Type:        p.Type,
		})
	}

	cc.DB.
		Preload("Menus").
		Preload("Products").
		First(&commande)

	c.JSON(http.StatusOK, gin.H{
		"message":  "Commande mise à jour avec succès",
		"commande": commande,
	})
}

// DeleteCommande supprime une commande
// @Summary Delete an existing commande
// @Description Delete an existing commande by ID
// @Tags commandes
// @Accept json
// @Produce json
// @Success 201 {object} models.Commande
// @Router /api/commandes/{id} [delete]
func (cc *CommandeController) DeleteCommande(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	if err := models.DeleteCommande(cc.DB, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression de la commande"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Commande supprimée avec succès"})
}
