package controllers

import (
	"LearningCampusKabre/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

func RefAuthController(db *gorm.DB) *AuthController {
	return &AuthController{DB: db}
}

// Register godoc
// @Summary Enregistrer un utilisateur
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.RegisterInput true "Données utilisateur"
// @Success 200 {object} map[string]string
// @Router /auth/register [post]
func (ac *AuthController) Register(c *gin.Context) {
	var input models.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides"})
		return
	}

	user := models.User{
		Email:    input.Email,
		Password: input.Password,
		Role:     models.UserRole(input.Role),
	}

	var count int64
	ac.DB.Model(&models.User{}).Where("email = ?", user.Email).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "L'email existe déjà"})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	if err := models.CreateUser(ac.DB, &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création de l'utilisateur"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Utilisateur enregistré avec succès"})
}

// Login godoc
// @Summary Connexion utilisateur
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body models.LoginInput true "Email et mot de passe"
// @Success 200 {object} map[string]string
// @Router /auth/login [post]
func (ac *AuthController) Login(c *gin.Context) {
	var input models.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données de connexion invalides"})
		return
	}

	var user models.User
	if err := ac.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email ou mot de passe incorrect"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email ou mot de passe incorrect"})
		return
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    string(user.Role),
		"exp":     time.Now().Add(2 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SIGNATURE_KEY")))

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
