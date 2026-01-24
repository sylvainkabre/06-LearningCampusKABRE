package main

import (
	"fmt"
	"log"
	"os"

	"LearningCampusKabre/models"
	"LearningCampusKabre/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	// Charger le fichier .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erreur lors du chargement du fichier .env")
	}

	// Récupérer les variables
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	sslmode := os.Getenv("DB_SSLMODE")
	signatureKey := os.Getenv("JWT_SIGNATURE_KEY")

	// Construire le DSN
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbname, port, sslmode)

	fmt.Println("DSN:", dsn)
	fmt.Println("Signature key:", signatureKey)

	// Connexion à la base de données
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Erreur de connexion à la base de données:", err)
	}

	// Migration automatique
	db.AutoMigrate(&models.Product{}, &models.User{})

	// Initialiser Gin
	router := gin.Default()

	// Middleware pour injecter la DB dans chaque requête
	router.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	// Configurer les routes
	routes.SetupProductRoutes(router, db)
	routes.SetupUserRoutes(router, db)

	// Démarrer le serveur
	log.Println("✅ Serveur démarré sur http://localhost:8080")
	router.Run(":8080")
}
