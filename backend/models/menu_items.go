package models

import (
	"time"

	"gorm.io/gorm"
)

type MenuItem struct {
	ID        uint `json:"id" gorm:"primaryKey"`
	MenuID    uint `json:"menu_id"`
	ProductID uint `json:"product_id"`

	Product Product `json:"product" gorm:"foreignKey:ProductID"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
