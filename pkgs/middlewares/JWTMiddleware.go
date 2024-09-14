package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lordofthemind/htmx_GO/internals/responses"
	"github.com/lordofthemind/htmx_GO/pkgs/tokens"
)

func AuthTokenMiddleware(tokenManager tokens.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("SuperUserAuthorization")
		if err != nil {
			response := responses.NewResponse(
				c,
				http.StatusUnauthorized,
				"Unauthorized",
				nil,
				"Failed to get token from cookie",
			)
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		}

		payload, err := tokenManager.ValidateToken(token)
		if err != nil {
			response := responses.NewResponse(
				c,
				http.StatusUnauthorized,
				"Invalid token",
				nil,
				err.Error(),
			)
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		}

		c.Set("userID", payload.ID)         // Use payload ID or other necessary field
		c.Set("username", payload.Username) // Optionally set username if needed
		c.Next()
	}
}
