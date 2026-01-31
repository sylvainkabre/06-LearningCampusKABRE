package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

// Définition du type pour le menu
type Menu struct {
	ID          uint                  `json:"id" gorm:"primaryKey"`
	Name        string                `json:"name" gorm:"not null"`
	Price       float32               `json:"price" gorm:"not null"`
	ImageURL    string                `json:"image_url"`
	Description string                `json:"description" gorm:"type:text"`
	MenuItems   []MenuItem            `json:"menu_items" gorm:"foreignKey:MenuID"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
	DeletedAt   soft_delete.DeletedAt `gorm:"softDelete:milli"`
}

// CreateMenu crée un nouveau menu
func CreateMenu(db *gorm.DB, menu *Menu) error {
	return db.Create(menu).Error
}

// GetAllMenus récupère tous les menus
func GetAllMenus(db *gorm.DB) ([]Menu, error) {
	var menus []Menu
	err := db.Find(&menus).Error
	return menus, err
}

// GetMenuByID récupère un menu par son ID
func GetMenuByID(db *gorm.DB, id uint) (*Menu, error) {
	var menu Menu
	err := db.First(&menu, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("Le menu avec l'Id %d n'a pas été trouvé", id)
		}
		return nil, fmt.Errorf("erreur lors de la récupération: %w", err)
	}
	return &menu, nil
}

// UpdateMenu met à jour un menu existant
func UpdateMenu(db *gorm.DB, menu *Menu) error {
	return db.Save(menu).Error
}

// DeleteMenu supprime un menu
func DeleteMenu(db *gorm.DB, id uint) error {
	return db.Delete(&Menu{}, id).Error
}
