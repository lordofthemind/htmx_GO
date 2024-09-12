package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/lordofthemind/htmx_GO/internals/responses"
)

// ResponseStrategyMiddleware injects the appropriate response strategy into the context.
func ResponseStrategyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		acceptHeader := c.GetHeader("Accept")

		switch {
		case acceptHeader == "text/html" || c.GetHeader("HX-Request") == "true":
			// For HTML or HTMX requests
			template := "default.html"
			if c.FullPath() == "/superuser/register" {
				template = "register_success.html"
			} else if c.FullPath() == "/superuser/login" {
				template = "login_response.html"
			}
			c.Set("responseStrategy", &responses.HTMLResponseStrategy{Template: template})
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
