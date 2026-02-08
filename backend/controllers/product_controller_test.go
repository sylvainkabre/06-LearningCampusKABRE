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
func setupProductTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.Product{})
	return db
}

// Configuration du routeur pour les tests
func setupProductRouter(pc *ProductController) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/products", pc.CreateProduct)
	return r
}

// Test de la création d'un produit
func TestCreateProduct(t *testing.T) {
	db := setupProductTestDB()
	pc := NewProductController(db)
	router := setupProductRouter(pc)

	// Données envoyées dans la requête (adaptées à TA struct Product)
	body := map[string]interface{}{
		"name":         "Salade César",
		"price":        9.99,
		"is_available": true,
		"image_url":    "https://example.com/images/salade-cesar.jpg",
		"description":  "Une délicieuse salade composée de laitue, poulet grillé, croûtons et parmesan.",
		"type":         "entree",
	}

	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var created models.Product
	json.Unmarshal(w.Body.Bytes(), &created)

	assert.Equal(t, "Salade César", created.Name)
	assert.Equal(t, "9.99", created.Price.String()) // Convertir decimal.Decimal en string pour la comparaison sinon Testify compare float et decimal.decimal
	assert.Equal(t, true, created.IsAvailable)
	assert.Equal(t, "https://example.com/images/salade-cesar.jpg", created.ImageURL)
	assert.Equal(t, "Une délicieuse salade composée de laitue, poulet grillé, croûtons et parmesan.", created.Description)
	assert.Equal(t, "entree", string(created.Type)) // idem pour TypeProduct
	assert.NotZero(t, created.ID)
}
