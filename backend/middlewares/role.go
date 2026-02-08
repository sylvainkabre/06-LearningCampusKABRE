package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireRole(allowedRoles ...string) gin.HandlerFunc { // Plus simple d'utiliser... sinon avec []string : middlewares.RequireRole([]string{"admin", "manager"})
	return func(c *gin.Context) {
		roleValue, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Rôle non trouvé dans le token"})
			return
		}

		userRole, ok := roleValue.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Format de rôle invalide"})
			return
		}

		// Vérification du rôle
		for _, allowed := range allowedRoles {
			if userRole == allowed {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Accès interdit"})
	}
}
