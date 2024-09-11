package responses

import "github.com/gin-gonic/gin"

type ResponseStrategy interface {
	Respond(c *gin.Context, data interface{}, status int)
}

type JSONResponseStrategy struct{}

func (r *JSONResponseStrategy) Respond(c *gin.Context, data interface{}, status int) {
	c.JSON(status, data)
}

type HTMLResponseStrategy struct {
	Template string
}

func (r *HTMLResponseStrategy) Respond(c *gin.Context, data interface{}, status int) {
	c.HTML(status, r.Template, data)
}
