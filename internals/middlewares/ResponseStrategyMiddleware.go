package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/lordofthemind/htmx_GO/internals/responses"
)

func ResponseStrategyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		acceptHeader := c.GetHeader("Accept")

		var strategy responses.ResponseStrategy
		if acceptHeader == "application/json" {
			strategy = &responses.JSONResponseStrategy{}
		} else {
			strategy = &responses.HTMLResponseStrategy{
				Template: "login_response.html", // or any appropriate template
			}
		}
		c.Set("responseStrategy", strategy)
		c.Next()
	}
}
