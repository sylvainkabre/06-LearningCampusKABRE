package models

// @Description Données nécessaires pour créer un utilisateur
type RegisterInput struct {
	Email    string `json:"email" example:"john@example.com"`
	Password string `json:"password" example:"motdepasse123"`
	Role     string `json:"role" example:"admin" enums:"admin,preparer,receiver"`
}

// @Description Données nécessaires pour se connecter
type LoginInput struct {
	Email    string `json:"email" example:"john@example.com"`
	Password string `json:"password" example:"motdepasse123"`
}
