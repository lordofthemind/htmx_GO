package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lordofthemind/htmx_GO/internals/tokens"
)

// JWTAuthMiddleware validates JWT tokens using the provided TokenManager.
func JWTAuthMiddleware(tokenManager tokens.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the JWT token from cookies
		token, err := c.Cookie("SuperUserAuthorization")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Validate the token using the TokenManager
		claims, err := tokenManager.ValidateJWT(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Set user ID in context
		c.Set("userID", claims["user_id"])
		c.Next()
	}
}
