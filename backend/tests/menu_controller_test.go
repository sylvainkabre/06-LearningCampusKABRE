package tests

import (
	"LearningCampusKabre/controllers"
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

// Base de données en mémoire
func setupMenuTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.Menu{}, &models.MenuItem{})
	return db
}

// Router pour les tests
func setupMenuRouter(pc *controllers.MenuController) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/menus", pc.CreateMenu)
	r.GET("/menus", pc.GetAllMenus)
	return r
}

/////////////////////////////////////
// TEST : CREATE MENU
/////////////////////////////////////

func TestCreateMenu(t *testing.T) {
	db := setupMenuTestDB()

	// On doit migrer Product et MenuItem aussi
	db.AutoMigrate(&models.Product{}, &models.MenuItem{})

	// On crée 2 produits disponibles
	p1 := models.Product{
		Name:        "Entrée",
		Price:       decimal.NewFromFloat(5.50),
		IsAvailable: true,
		Type:        models.TypeProduct("entree"),
	}
	p2 := models.Product{
		Name:        "Plat",
		Price:       decimal.NewFromFloat(16.00),
		IsAvailable: true,
		Type:        models.TypeProduct("plat"),
	}
	db.Create(&p1)
	db.Create(&p2)

	pc := controllers.RefMenuController(db)
	router := setupMenuRouter(pc)

	body := map[string]interface{}{
		"name":        "Menu estival",
		"price":       "21.50",
		"items":       []int{int(p1.ID), int(p2.ID)}, // IMPORTANT
		"image_url":   "https://example.com/images/menu_estival.jpg",
		"description": "Menu aux saveurs estivales",
	}

	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/menus", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// On lit la réponse complète
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	// On récupère le menu dans "menu"
	menuMap := response["menu"].(map[string]interface{})

	assert.Equal(t, "Menu estival", menuMap["name"])
	assert.Equal(t, "21.5", menuMap["price"])
	assert.Equal(t, "https://example.com/images/menu_estival.jpg", menuMap["image_url"])
	assert.Equal(t, "Menu aux saveurs estivales", menuMap["description"])
	assert.NotZero(t, menuMap["id"])
}

func TestGetAllMenus(t *testing.T) {
	db := setupMenuTestDB()

	// Migration complète
	db.AutoMigrate(&models.Product{}, &models.MenuItem{})

	// On crée 2 produits disponibles
	p1 := models.Product{
		Name:        "Entrée",
		Price:       decimal.NewFromFloat(5.50),
		IsAvailable: true,
		Type:        models.TypeProduct("entree"),
	}
	p2 := models.Product{
		Name:        "Plat",
		Price:       decimal.NewFromFloat(16.00),
		IsAvailable: true,
		Type:        models.TypeProduct("plat"),
	}
	db.Create(&p1)
	db.Create(&p2)

	// On crée un menu
	menu := models.Menu{
		Name:        "Menu estival",
		Description: "Menu aux saveurs estivales",
		Price:       decimal.NewFromFloat(21.50),
		ImageURL:    "https://example.com/images/menu_estival.jpg",
	}
	db.Create(&menu)

	// On crée les items associés
	db.Create(&models.MenuItem{
		MenuID:      menu.ID,
		Name:        p1.Name,
		Price:       p1.Price,
		ImageURL:    p1.ImageURL,
		Description: p1.Description,
		Type:        p1.Type,
	})
	db.Create(&models.MenuItem{
		MenuID:      menu.ID,
		Name:        p2.Name,
		Price:       p2.Price,
		ImageURL:    p2.ImageURL,
		Description: p2.Description,
		Type:        p2.Type,
	})

	pc := controllers.RefMenuController(db)
	router := setupMenuRouter(pc)

	// Requête GET
	req, _ := http.NewRequest("GET", "/menus", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Lecture du JSON
	var menus []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &menus)

	// On doit avoir 1 menu
	assert.Len(t, menus, 1)

	m := menus[0]

	assert.Equal(t, "Menu estival", m["name"])
	assert.Equal(t, "21.5", m["price"])
	assert.Equal(t, "https://example.com/images/menu_estival.jpg", m["image_url"])
	assert.Equal(t, "Menu aux saveurs estivales", m["description"])
	assert.NotZero(t, m["id"])

	// Vérification des items préloadés
	items := m["menu_items"].([]interface{})
	assert.Len(t, items, 2)
}
