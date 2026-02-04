package models

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

// Définition du type pour le produit
type StatusType string

// Constantes pour les valeurs possibles
const (
	StatusPending   StatusType = "pending"
	StatusPreparing StatusType = "preparing"
	StatusReady     StatusType = "ready"
	StatusDelivered StatusType = "delivered"
)

// Une commande peut être composée de plusieurs menus et produits
type Commande struct {
	ID        uint                  `json:"id" gorm:"primaryKey"`
	Menus     []CommandeMenu        `json:"menus" gorm:"foreignKey:CommandeID"`
	Products  []CommandeProduct     `json:"products" gorm:"foreignKey:CommandeID"`
	Status    StatusType            `json:"status" gorm:"not null"`
	Price     decimal.Decimal       `json:"price" gorm:"type:decimal(10,2);not null"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
	DeletedAt soft_delete.DeletedAt `gorm:"softDelete:milli" swaggerignore:"true"`
}

// CreateCommande crée une nouvelle commande
func CreateCommande(db *gorm.DB, commande *Commande) error {
	return db.Create(commande).Error
}

// GetAllComm récupère toutes les commandes
func GetAllComm(db *gorm.DB) ([]Commande, error) {
	var commandes []Commande
	err := db.Find(&commandes).Error
	return commandes, err
}

func GetCommandeById(db *gorm.DB, id uint) (*Commande, error) {
	var commande Commande
	err := db.First(&commande, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("La commande avec l'Id %d n'a pas été trouvée", id)
		}
		return nil, fmt.Errorf("erreur lors de la récupération: %w", err)
	}
	return &commande, nil
}

// UpdateCommande met à jour une commande existante
func UpdateCommande(db *gorm.DB, commande *Commande) error {
	return db.Save(commande).Error
}

// DeleteCommande supprime une commande
func DeleteCommande(db *gorm.DB, id uint) error {
	return db.Delete(&Commande{}, id).Error
}
