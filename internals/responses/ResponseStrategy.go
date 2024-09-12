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
	Template string // default template
}

// Respond sends an HTML response using a specified template.
func (r *HTMLResponseStrategy) Respond(c *gin.Context, data interface{}, status int) {
	// If the data includes a "template" field, use it as the template name
	if dataMap, ok := data.(map[string]interface{}); ok {
		if templateName, exists := dataMap["template"]; exists {
			if templateNameStr, isString := templateName.(string); isString {
				c.HTML(status, templateNameStr, dataMap)
				return
			}
		}
	}

	// Fall back to the default template if no "template" is provided
	c.HTML(status, r.Template, data)
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
