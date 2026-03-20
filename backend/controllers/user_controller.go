package controllers

import (
	"LearningCampusKabre/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

func RefUserController(db *gorm.DB) *UserController {
	return &UserController{DB: db}
}

// CreateUser godoc
// @Summary Créer un utilisateur
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "Utilisateur à créer"
// @Success 201 {object} models.User
// @Router /users [post]
// @Security BearerAuth
func (uc *UserController) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !user.Role.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Rôle utilisateur invalide"})
		return
	}

	if err := models.CreateUser(uc.DB, &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création de l'utilisateur"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// GetAllUsers godoc
// @Summary Liste des utilisateurs
// @Tags users
// @Produce json
// @Success 200 {array} models.User
// @Router /users [get]
// @Security BearerAuth
func (uc *UserController) GetAllUsers(c *gin.Context) {
	var users []models.User
	if err := uc.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des utilisateurs"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetUserByID godoc
// @Summary Récupérer un utilisateur
// @Tags users
// @Produce json
// @Param id path int true "ID utilisateur"
// @Success 200 {object} models.User
// @Router /users/{id} [get]
// @Security BearerAuth
func (uc *UserController) GetUserByID(c *gin.Context) {
	id := c.Param("id")

	var userID uint
	if _, err := fmt.Sscan(id, &userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	user, err := models.GetUserByID(uc.DB, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser godoc
// @Summary Mettre à jour un utilisateur
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID utilisateur"
// @Param user body models.User true "Données utilisateur"
// @Success 200 {object} models.User
// @Router /users/{id} [put]
// @Security BearerAuth
func (uc *UserController) UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var userID uint
	if _, err := fmt.Sscan(id, &userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	existingUser, err := models.GetUserByID(uc.DB, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Utilisateur non trouvé"})
		return
	}

	var updatedUser models.User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !updatedUser.Role.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Rôle utilisateur invalide"})
		return
	}

	updatedUser.ID = existingUser.ID

	if err := models.UpdateUser(uc.DB, &updatedUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour de l'utilisateur"})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

// DeleteUser godoc
// @Summary Supprimer un utilisateur
// @Tags users
// @Param id path int true "ID utilisateur"
// @Success 200 {object} map[string]string
// @Router /users/{id} [delete]
// @Security BearerAuth
func (ctrl *UserController) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := ctrl.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Utilisateur non trouvé"})
		return
	}

	if err := ctrl.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du soft delete"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Utilisateur soft supprimé"})
}
