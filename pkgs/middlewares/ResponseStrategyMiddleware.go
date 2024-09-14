package middlewares

import (
	"log"

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
			log.Println("Response Strategy set to HTML Response")
		case acceptHeader == "application/json":
			// For JSON requests
			c.Set("responseStrategy", &responses.JSONResponseStrategy{})
			log.Println("Response Strategy set to JSON Response")
		default:
			// Fallback to HTMLResponseStrategy for any other cases
			c.Set("responseStrategy", &responses.HTMLResponseStrategy{Template: "default.html"})
			log.Println("Response Strategy set to HTML Response in default")
		}

		c.Next()
	}
}
