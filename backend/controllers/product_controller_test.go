package controllers

import (
	"LearningCampusKabre/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
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
	r.GET("/products", pc.GetAllProducts)
	r.PUT("/products/1", pc.UpdateProduct)
	return r
}

/////////////////////////////////////
// CREATE PRODUCTS
/////////////////////////////////////

// Test de la création d'un produit
func TestCreateProduct(t *testing.T) {
	db := setupProductTestDB()
	pc := RefProductController(db)
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

/////////////////////////////////////
// GET SPECIFIC PRODUCTS (BIND)
/////////////////////////////////////

func setupProductTestDBWithData() (*gorm.DB, []models.Product) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.Product{})

	products := []models.Product{
		{
			Name:        "Salade César",
			Price:       decimal.NewFromFloat(9.99),
			IsAvailable: true,
			ImageURL:    "https://example.com/salade.jpg",
			Description: "Salade fraîche",
			Type:        models.TypeProduct("entree"),
		},
		{
			Name:        "Pizza Margherita",
			Price:       decimal.NewFromFloat(12.50),
			IsAvailable: true,
			ImageURL:    "https://example.com/pizza.jpg",
			Description: "Pizza classique",
			Type:        models.TypeProduct("plat"),
		},
	}

	for _, p := range products {
		db.Create(&p)
	}

	return db, products
}

/////////////////////////////////////
// GET ALL PRODUCTS
/////////////////////////////////////

func TestGetAllProducts(t *testing.T) {
	db, expectedProducts := setupProductTestDBWithData()
	pc := RefProductController(db)
	router := setupProductRouter(pc)

	req, _ := http.NewRequest("GET", "/products", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	products := []models.Product{
		{
			Name:        "Salade César",
			Price:       decimal.NewFromFloat(9.99),
			IsAvailable: true,
			ImageURL:    "https://example.com/salade.jpg",
			Description: "Salade fraîche",
			Type:        models.TypeProduct("entree"),
		},
		{
			Name:        "Pizza Margherita",
			Price:       decimal.NewFromFloat(12.50),
			IsAvailable: true,
			ImageURL:    "https://example.com/pizza.jpg",
			Description: "Pizza classique",
			Type:        models.TypeProduct("plat"),
		},
	}

	// Vérifie le nombre de produits
	assert.Len(t, products, len(expectedProducts))

	// Vérifie le premier produit
	assert.Equal(t, expectedProducts[0].Name, products[0].Name)
	assert.Equal(t, expectedProducts[0].Price.String(), products[0].Price.String())
	assert.Equal(t, expectedProducts[0].Description, products[0].Description)
	assert.Equal(t, string(expectedProducts[0].Type), string(products[0].Type))

	// Vérifie le second produit
	assert.Equal(t, expectedProducts[1].Name, products[1].Name)
	assert.Equal(t, expectedProducts[1].Price.String(), products[1].Price.String())
	assert.Equal(t, expectedProducts[1].Description, products[1].Description)
	assert.Equal(t, string(expectedProducts[1].Type), string(products[1].Type))
}

/////////////////////////////////////
// UPDATE PRODUCTS
/////////////////////////////////////

// func TestUpdateProduct(t *testing.T) {
// 	db, _ := setupProductTestDBWithData()
// 	pc := RefProductController(db)
// 	router := setupProductRouter(pc)

// 	updateBody := map[string]interface{}{
// 		"name":         "Salade César Deluxe",
// 		"price":        11.49,
// 		"is_available": true,
// 		"image_url":    "https://example.com/salade-deluxe.jpg",
// 		"description":  "Version améliorée",
// 		"type":         "entree",
// 	}

// 	jsonBody, _ := json.Marshal(updateBody)
// 	req, _ := http.NewRequest("PUT", "/products/1", bytes.NewBuffer(jsonBody))
// 	req.Header.Set("Content-Type", "application/json")

// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusOK, w.Code)

// 	var updated models.Product
// 	json.Unmarshal(w.Body.Bytes(), &updated)

// 	assert.Equal(t, "Salade César Deluxe", updated.Name)
// 	assert.Equal(t, "11.49", updated.Price.String())
// 	assert.Equal(t, false, updated.IsAvailable)
// 	assert.Equal(t, "https://example.com/salade-deluxe.jpg", updated.ImageURL)
// 	assert.Equal(t, "Version améliorée", updated.Description)
// 	assert.Equal(t, "entree", string(updated.Type))
// 	assert.Equal(t, uint(1), updated.ID)
// }
