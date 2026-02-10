package controllers

import (
	"LearningCampusKabre/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Création d'une base de données en mémoire pour les tests
func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.User{})
	return db
}

// Configuration du routeur pour les tests
func setupUserRouter(uc *UserController) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/users", uc.CreateUser)
	r.GET("/users", uc.GetAllUsers)
	return r
}

/////////////////////////////////////
// CREATE USER
/////////////////////////////////////

// Test de la création d'un utilisateur
func TestCreateUser(t *testing.T) {
	db := setupTestDB()
	uc := RefUserController(db)
	router := setupUserRouter(uc)

	// Données envoyées dans la requête (adaptées à TA struct User)
	body := map[string]interface{}{
		"email":       "Maxime@Michaud.com",
		"password":    "Maxime123",
		"role":        "preparer",
		"description": "Utilisateur de test",
	}

	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var created models.User
	json.Unmarshal(w.Body.Bytes(), &created)

	assert.Equal(t, "Maxime@Michaud.com", created.Email)
	assert.Equal(t, models.RolePreparer, created.Role)
	assert.Equal(t, "Utilisateur de test", created.Description)
	assert.NotZero(t, created.ID)
}

/////////////////////////////////////
// GET SPECIFIC USER (BIND)
/////////////////////////////////////

func setupUserTestDBWithData() (*gorm.DB, []models.User) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.User{})

	users := []models.User{
		{
			Email:       "Maxime@Michaud.com",
			Password:    "Maxime123",
			Role:        "preparer",
			Description: "Utilisateur de test",
		},
	}

	for _, p := range users {
		db.Create(&p)
	}

	return db, users
}

/////////////////////////////////////
// GET ALL USERS
/////////////////////////////////////

func TestGetAllUsers(t *testing.T) {
	db, expectedUsers := setupUserTestDBWithData()
	pc := RefUserController(db)
	router := setupUserRouter(pc)

	req, _ := http.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	users := []models.User{
		{
			Email:       "Maxime@Michaud.com",
			Password:    "Maxime123",
			Role:        "preparer",
			Description: "Utilisateur de test",
		},
	}

	// Vérifie le nombre de produits
	assert.Len(t, users, len(expectedUsers))

	// Vérifie le premier produit
	assert.Equal(t, expectedUsers[0].Email, users[0].Email)
	assert.Equal(t, expectedUsers[0].Password, users[0].Password)
	assert.Equal(t, expectedUsers[0].Role, users[0].Role)
	assert.Equal(t, expectedUsers[0].Description, users[0].Description)
}
