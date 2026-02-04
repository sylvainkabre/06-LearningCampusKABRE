package models

import (
	"time"

	"github.com/shopspring/decimal"
)

// CommandeProduct représente un produit dans une commande
type CommandeProduct struct {
	ID          uint            `json:"id" gorm:"primaryKey"`
	CommandeID  uint            `json:"commande_id"`
	ProductID   uint            `json:"product_id"`
	Name        string          `json:"name"`
	Price       decimal.Decimal `json:"price"`
	ImageURL    string          `json:"image_url"`
	Description string          `json:"description"`
	Type        TypeProduct     `json:"type"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}
