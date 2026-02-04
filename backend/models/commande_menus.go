package models

import (
	"time"

	"github.com/shopspring/decimal"
)

// CommandeMenu représente un menu dans une commande
type CommandeMenu struct {
	ID          uint            `json:"id" gorm:"primaryKey"`
	CommandeID  uint            `json:"commande_id"`
	MenuID      uint            `json:"menu_id"`
	Name        string          `json:"name"`
	Price       decimal.Decimal `json:"price"`
	Description string          `json:"description"`
	ImageURL    string          `json:"image_url"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}
