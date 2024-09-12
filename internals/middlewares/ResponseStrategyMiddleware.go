package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/lordofthemind/htmx_GO/internals/responses"
)

func ResponseStrategyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		acceptHeader := c.GetHeader("Accept")

		switch {
		case acceptHeader == "text/html" || c.GetHeader("HX-Request") == "true":
			// For HTML or HTMX requests, no template specified here
			c.Set("responseStrategy", &responses.HTMLResponseStrategy{})
		case acceptHeader == "application/json":
			// For JSON requests
			c.Set("responseStrategy", &responses.JSONResponseStrategy{})
		default:
			// Fallback to HTMLResponseStrategy for any other cases
			c.Set("responseStrategy", &responses.HTMLResponseStrategy{Template: "default.html"})
		}

		c.Next()
	}
}
