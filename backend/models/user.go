package models

//Création de CreateUser, GetUserByID, UpdateUser, DeleteUser

import (
	"time"

	"fmt"

	"gorm.io/gorm"

	"gorm.io/plugin/soft_delete"
)

// Définition du type pour le rôle
type UserRole string

// Constantes pour les valeurs possibles
const (
	RoleAdmin     UserRole = "admin"
	RoleModerator UserRole = "preparer"
	RoleUser      UserRole = "receiver"
)

// Méthode pour valider si un rôle est valide
func (r UserRole) IsValid() bool {
	switch r {
	case RoleAdmin, RoleModerator, RoleUser:
		return true
	}
	return false
}

// Structure User avec le type UserRole
type User struct {
	ID          uint                  `json:"id" gorm:"primaryKey"`
	Email       string                `json:"email" gorm:"not null"`
	Role        UserRole              `json:"role" gorm:"type:varchar(20);not null"`
	Description string                `json:"description" gorm:"type:text"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
	DeletedAt   soft_delete.DeletedAt `gorm:"softDelete:milli"`
	Password    string                `json:"password" binding:"required,min=6"`
}

// Fonction pour créer un utilisateur en base de données
func CreateUser(db *gorm.DB, user *User) error {
	result := db.Create(user)
	return result.Error
}

// Récupérer un utilisateur par ID
func GetUserByID(db *gorm.DB, id uint) (*User, error) {
	var user User
	err := db.First(&user, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("L'utilisateur avec l'Id %d n'a pas été trouvé", id)
		}
		return nil, fmt.Errorf("erreur lors de la récupération: %w", err)
	}

	return &user, nil
}

// Mettre à jour un utilisateur
func UpdateUser(db *gorm.DB, user *User) error {
	return db.Save(user).Error
}

// Supprimer un utilisateur
func DeleteUser(db *gorm.DB, id uint) error {
	return db.Delete(&User{}, id).Error
}

// LoginUser vérifie les informations d'identification de l'utilisateur
func LoginUser(db *gorm.DB, name string, password string) (*User, error) {
	var user User
	err := db.Where("name = ? AND password = ?", name, password).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("Nom d'utilisateur ou mot de passe incorrect")
		}
		return nil, fmt.Errorf("erreur lors de la récupération: %w", err)
	}

	return &user, nil
}
