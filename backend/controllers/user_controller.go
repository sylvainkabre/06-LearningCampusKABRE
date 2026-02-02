package controllers

//Appel de CreateUser, GetUserByID, UpdateUser, DeleteUser

import (
	"LearningCampusKabre/models"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type CustomClaims struct {
	UserID uint
	Role   string
	jwt.RegisteredClaims
}

type UserController struct {
	DB *gorm.DB
}

func NewUserController(db *gorm.DB) *UserController {
	return &UserController{DB: db}
}

// CreateUser godoc
// @Summary Créer un utilisateur
// @Description Ajoute un nouvel utilisateur en base
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "Utilisateur à créer"
// @Success 201 {object} models.User
// @Router /users [post]
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
// @Description Retourne tous les utilisateurs
// @Tags users
// @Produce json
// @Success 200 {array} models.User
// @Router /users [get]
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
// @Description Retourne un utilisateur selon son ID
// @Tags users
// @Produce json
// @Param id path int true "ID utilisateur"
// @Success 200 {object} models.User
// @Router /users/{id} [get]
func (uc *UserController) GetUserByID(c *gin.Context) {
	id := c.Param("id")

	// Convertir l'ID string en uint
	var userID uint
	if _, err := fmt.Sscan(id, &userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	// Récupérer l'utilisateur depuis la base de données
	user, err := models.GetUserByID(uc.DB, userID)
	if err != nil {
		// ✅ Utiliser err.Error() au lieu d'un message codé en dur
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser godoc
// @Summary Mettre à jour un utilisateur
// @Description Modifie les informations d’un utilisateur existant
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID utilisateur"
// @Param user body models.User true "Données utilisateur"
// @Success 200 {object} models.User
// @Router /users/{id} [put]
func (uc *UserController) UpdateUser(c *gin.Context) {
	id := c.Param("id")

	// Convertir l'ID string en uint
	var userID uint
	if _, err := fmt.Sscan(id, &userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	// Vérifier que l'utilisateur existe
	existingUser, err := models.GetUserByID(uc.DB, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Utilisateur non trouvé"})
		return
	}

	// Parser les nouvelles données
	var updatedUser models.User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Valider le rôle si fourni
	if !updatedUser.Role.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Rôle utilisateur invalide"})
		return
	}

	// Conserver l'ID original
	updatedUser.ID = existingUser.ID

	// Mettre à jour en base de données
	if err := models.UpdateUser(uc.DB, &updatedUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour de l'utilisateur"})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

// DeleteUser godoc
// @Summary Supprimer un utilisateur (soft delete)
// @Description Marque un utilisateur comme supprimé
// @Tags users
// @Param id path int true "ID utilisateur"
// @Success 200 {object} map[string]string
// @Router /users/{id} [delete]
func (ctrl *UserController) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User

	// Vérifier que l'utilisateur existe (non supprimé)
	if err := ctrl.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Utilisateur non trouvé"})
		return
	}

	// Soft delete automatique (remplit deleted_at)
	if err := ctrl.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du soft delete"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Utilisateur soft supprimé"})
}

// Register godoc
// @Summary Enregistrer un utilisateur
// @Description Crée un utilisateur avec mot de passe hashé
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.User true "Données utilisateur"
// @Success 200 {object} map[string]string
// @Router /auth/register [post]
func Register(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil { //ShouldBindJSON s'assure que les données JSON reçues sont correctement liées à la structure User
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides"})
		return
	}

	if !user.Role.IsValid() { // IsValid vérifie que le rôle de l'utilisateur est valide
		c.JSON(http.StatusBadRequest, gin.H{"error": "Rôle utilisateur invalide"})
		return
	}

	// Vérifier si l'email existe déjà
	var count int64
	db.Model(&models.User{}).Where("email = ?", user.Email).Count(&count) // Model spécifie la table, Where ajoute une condition, Count compte les enregistrements correspondants
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "L'email existe déjà"})
		return
	}

	// Hashage du mot de passe pour la sécurité en base de données
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost) // GenerateFromPassword crée un hachage sécurisé du mot de passe
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur serveur"})
		return
	}
	user.Password = string(hashedPassword)

	if err := models.CreateUser(db, &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création de l'utilisateur"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Utilisateur enregistré avec succès"})
}

// Login godoc
// @Summary Connexion utilisateur
// @Description Retourne un token JWT si les identifiants sont valides
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body map[string]string true "Email et mot de passe"
// @Success 200 {object} map[string]string
// @Router /auth/login [post]
func Login(c *gin.Context) {
	var loginData struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données de connexion invalides"})
		return
	}
	db := c.MustGet("db").(*gorm.DB)
	var user models.User
	if err := db.Where("email = ?", loginData.Email).First(&user).Error; err != nil { // On vérifie uniquement email dans la base de données car le mot de passe est haché au moment du register
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email ou mot de passe incorrect 1"})
		return
	}
	// Vérifier le mot de passe haché
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email ou mot de passe incorrect 2"})
		return
	}

	claim := &CustomClaims{
		UserID: user.ID,
		Role:   string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SIGNATURE_KEY"))) // A mettre dans fichier .env pour plus de sécurité

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la génération du token"})
		return
	}

	c.JSON(http.StatusOK, "Token :"+tokenString)
}
