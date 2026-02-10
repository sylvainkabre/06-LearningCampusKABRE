package models

//Création de CreateProduct, GetAllProducts, GetProductByID, UpdateProduct, DeleteProduct

import (
	"time"

	"fmt"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"gorm.io/plugin/soft_delete"
)

// Définition du type pour le produit
type TypeProduct string

// Constantes pour les valeurs possibles
const (
	TypeEntree  TypeProduct = "entree"
	TypePlat    TypeProduct = "plat"
	TypeDessert TypeProduct = "dessert"
	TypeBoisson TypeProduct = "boisson"
)

// Méthode pour valider si un type de produit est valide
func (r TypeProduct) IsValid() bool {
	switch r {
	case TypeEntree, TypePlat, TypeDessert:
		return true
	}
	return false
}

type Product struct {
	ID          uint                  `json:"id" gorm:"primaryKey" example:"1"`
	Name        string                `json:"name" gorm:"not null" example:"Salade César"`
	Price       decimal.Decimal       `json:"price" gorm:"not null" example:"9.99"`
	IsAvailable bool                  `json:"is_available" gorm:"default:true" example:"true"`
	ImageURL    string                `json:"image_url" gorm:"type:text" example:"https://example.com/images/salade-cesar.jpg"`
	Description string                `json:"description" gorm:"type:text" example:"Une délicieuse salade composée de laitue, poulet grillé, croûtons et parmesan."`
	Type        TypeProduct           `json:"type" gorm:"not null" example:"entree"`
	CreatedAt   time.Time             `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time             `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   soft_delete.DeletedAt `gorm:"softDelete:milli" json:"-"`
}

// CreateProduct crée un nouveau produit
func CreateProduct(db *gorm.DB, product *Product) error {
	return db.Create(product).Error
}

// GetAllProducts récupère tous les produits
func GetAllProducts(db *gorm.DB) ([]Product, error) {
	var products []Product
	err := db.Find(&products).Error
	return products, err
}

// GetProductByID récupère un produit par son ID
func GetProductByID(db *gorm.DB, id uint) (*Product, error) {
	var product Product
	err := db.First(&product, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("Le produit avec l'Id %d n'a pas été trouvé", id)
		}
		return nil, fmt.Errorf("erreur lors de la récupération: %w", err)
	}
	return &product, nil
}

// UpdateProduct met à jour un produit existant
func UpdateProduct(db *gorm.DB, product *Product) error {
	return db.Save(product).Error
}

// DeleteProduct supprime un produit
func DeleteProduct(db *gorm.DB, id uint) error {
	return db.Delete(&Product{}, id).Error
}

func GetAvailableItemsByType(db *gorm.DB, Type string) ([]Product, error) {
	var items []Product
	err := db.Where("type = ? AND is_available = ?", Type, true).Find(&items).Error
	return items, err
}
