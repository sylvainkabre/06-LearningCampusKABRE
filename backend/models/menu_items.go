package models

import (
	"time"

	"gorm.io/gorm"
)

// MenuItem represents a menu item entity
type MenuItem struct {
	ID          uint        `json:"id" gorm:"primaryKey"`
	MenuID      uint        `json:"menu_id"`
	Name        string      `json:"name"`
	Price       float32     `json:"price"`
	ImageURL    string      `json:"image_url"`
	Description string      `json:"description"`
	Type        TypeProduct `json:"type"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index" swaggerignore:"true"`
}
