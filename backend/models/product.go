package models

//Création de CreateProduct, GetAllProducts, GetProductByID, UpdateProduct, DeleteProduct

import (
	"time"

	"fmt"

	"gorm.io/gorm"

	"gorm.io/plugin/soft_delete"
)

type Product struct {
	ID          uint                  `json:"id" gorm:"primaryKey"`
	Name        string                `json:"name" gorm:"not null"`
	Price       float32               `json:"price" gorm:"not null"`
	IsAvailable bool                  `json:"available" gorm:"default:true"`
	ImageURL    string                `json:"image_url"`
	Description string                `json:"description" gorm:"type:text"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
	DeletedAt   soft_delete.DeletedAt `gorm:"softDelete:milli"`
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
