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
func setupMenuTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.Menu{})
	return db
}

// Configuration du routeur pour les tests
func setupMenuRouter(pc *MenuController) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/menus", pc.CreateMenu)
	r.GET("/menus", pc.GetAllMenus)
	return r
}

/////////////////////////////////////
// CREATE MENU
/////////////////////////////////////

// Test de la création d'un produit
func TestCreateMenu(t *testing.T) {
	db := setupMenuTestDB()
	pc := RefMenuController(db)
	router := setupMenuRouter(pc)

	// Données envoyées dans la requête (adaptées à TA struct Product)
	body := map[string]interface{}{
		"name":        "Menu estival",
		"price":       21.50,
		"menu_items":  []int{1, 2},
		"image_url":   "https://example.com/images/menu_estival.jpg",
		"description": "Menu aux saveurs estivales",
	}

	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/menus", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var created models.Menu
	json.Unmarshal(w.Body.Bytes(), &created)

	assert.Equal(t, "Menu estival", created.Name)
	assert.Equal(t, "21.50", created.Price.String()) // Convertir decimal.Decimal en string pour la comparaison sinon Testify compare float et decimal.decimal
	assert.Equal(t, "https://example.com/images/menu_estival.jpg", created.ImageURL)
	assert.Equal(t, "Menu aux saveurs estivales", created.Description)
	assert.NotZero(t, created.ID)
}
