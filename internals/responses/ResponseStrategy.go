package responses

import (
	"log"

	"github.com/gin-gonic/gin"
)

// ResponseStrategy defines the interface for responding with different formats.
type ResponseStrategy interface {
	Respond(c *gin.Context, data interface{}, status int)
}

// JSONResponseStrategy implements ResponseStrategy for JSON responses.
type JSONResponseStrategy struct{}

// Respond sends a JSON response.
func (r *JSONResponseStrategy) Respond(c *gin.Context, data interface{}, status int) {
	c.JSON(status, data)
}

// HTMLResponseStrategy implements ResponseStrategy for HTML responses.
type HTMLResponseStrategy struct {
	Template string
}

// Respond sends an HTML response using a specified template.
func (r *HTMLResponseStrategy) Respond(c *gin.Context, data interface{}, status int) {
	// Call c.HTML without capturing a return value, as it returns nothing
	c.HTML(status, r.Template, data)
	// Log if you want to track rendering or add additional error handling in Gin context
}

// DefaultResponseStrategy is a fallback strategy if none is set in the context.
type DefaultResponseStrategy struct{}

// Respond sends a plain text response as a default fallback.
func (r *DefaultResponseStrategy) Respond(c *gin.Context, data interface{}, status int) {
	c.String(status, "Default response - no strategy set")
}

// Retrieve the strategy with error handling.
func GetResponseStrategy(c *gin.Context) ResponseStrategy {
	strategy, exists := c.Get("responseStrategy")
	if !exists {
		log.Println("No response strategy set, using default")
		return &DefaultResponseStrategy{}
	}
	if responseStrategy, ok := strategy.(ResponseStrategy); ok {
		return responseStrategy
	}
	log.Println("Invalid response strategy set, using default")
	return &DefaultResponseStrategy{}
}
