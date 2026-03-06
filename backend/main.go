package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"LearningCampusKabre/models"
	"LearningCampusKabre/routes"

	_ "LearningCampusKabre/docs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title Learning Campus API
// @version 1.0
// @description API du projet LearningCampusKabre
// @host localhost:8080
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @security BearerAuth
func main() {

	// Charger .env uniquement en local
	_ = godotenv.Load()

	// Vérifier si on est en production Railway
	dsn := os.Getenv("DATABASE_URL")

	if dsn != "" {
		// --- MODE PRODUCTION ---
		fmt.Println("🌍 Mode production détecté")

		if !strings.Contains(dsn, "sslmode") {
			dsn += "?sslmode=require"
		}

	} else {
		// --- MODE LOCAL ---
		fmt.Println("💻 Mode local détecté")

		host := os.Getenv("DB_HOST")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		dbname := os.Getenv("DB_NAME")
		port := os.Getenv("DB_PORT")
		sslmode := os.Getenv("DB_SSLMODE")

		if sslmode == "" {
			sslmode = "disable"
		}

		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			host, user, password, dbname, port, sslmode,
		)
	}

	fmt.Println("DSN utilisé :", dsn)

	// Connexion DB
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Erreur de connexion à la base :", err)
	}

	// Migrations
	db.AutoMigrate(
		&models.Product{},
		&models.User{},
		&models.Menu{},
		&models.MenuItem{},
		&models.Commande{},
	)

	// Gin
	router := gin.Default()

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Injecter DB dans le contexte
	router.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	// Routes
	routes.SetupProductRoutes(router, db)
	routes.SetupUserRoutes(router, db)
	routes.SetupMenuRoutes(router, db)
	routes.SetupCommandesRoutes(router, db)

	// Port Railway ou 8080 en local
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("🚀 Serveur démarré sur http://localhost:" + port)
	router.Run(":" + port)
}
