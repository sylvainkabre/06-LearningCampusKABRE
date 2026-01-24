package middlewares

import (
	"net/http"

	"fmt"

	"github.com/gin-gonic/gin"
)

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		fmt.Println("role value in AdminOnly middleware:", role)
		// Si la clé n'existe pas ou n'est pas un booléen true → accès refusé
		if !exists || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Accès réservé aux administrateurs",
			})
			c.Abort()
			return
		}

		// Sinon on continue
		fmt.Println("AdminOnly exécuté")
		c.Next()
	}
}
