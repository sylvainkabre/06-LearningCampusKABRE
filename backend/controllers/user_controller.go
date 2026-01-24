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
	jwt.RegisteredClaims
}

type UserController struct {
	DB *gorm.DB
}

func NewUserController(db *gorm.DB) *UserController {
	return &UserController{DB: db}
}

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

func (uc *UserController) GetAllUsers(c *gin.Context) {
	var users []models.User
	if err := uc.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des utilisateurs"})
		return
	}
	c.JSON(http.StatusOK, users)
}

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

func (uc *UserController) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	// Convertir l'ID string en uint
	var userID uint
	if _, err := fmt.Sscan(id, &userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	// Vérifier que l'utilisateur existe
	_, err := models.GetUserByID(uc.DB, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Utilisateur non trouvé"})
		return
	}

	// Supprimer l'utilisateur
	if err := models.DeleteUser(uc.DB, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression de l'utilisateur"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Utilisateur supprimé avec succès"})
}

// Enregistrement et la connexion des utilisateurs approche similaire au cours

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

	// Vérifier si le nom existe déjà
	var count int64
	db.Model(&models.User{}).Where("name = ?", user.Name).Count(&count) // Model spécifie la table, Where ajoute une condition, Count compte les enregistrements correspondants
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Le nom d'utilisateur existe déjà"})
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

// LoginUser gère la connexion des utilisateurs

func Login(c *gin.Context) {
	var loginData struct {
		Name     string `json:"name" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données de connexion invalides"})
		return
	}
	db := c.MustGet("db").(*gorm.DB)
	var user models.User
	if err := db.Where("name = ?", loginData.Name).First(&user).Error; err != nil { // On vérifie uniquement name dans la base de données car le mot de passe est haché au moment du register
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Nom d'utilisateur ou mot de passe incorrect 1"})
		return
	}
	// Vérifier le mot de passe haché
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Nom d'utilisateur ou mot de passe incorrect 2"})
		return
	}

	claim := &CustomClaims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)), // Expiration du token dans 2 heures
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SIGNATURE_KEY"))) // A mettre dans fichier .env pour plus de sécurité

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la génération du token"})
		return
	}

	c.JSON(http.StatusOK, tokenString)
}
