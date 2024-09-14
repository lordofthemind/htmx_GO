package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lordofthemind/htmx_GO/pkgs/tokens"
)

func AuthTokenMiddleware(tokenManager tokens.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("SuperUserAuthorization")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		claims, err := tokenManager.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("userID", claims["user_id"])
		c.Next()
	}
}
