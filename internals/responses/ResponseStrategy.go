package responses

import (
	"fmt"
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
	// Log the data for debugging
	fmt.Printf("Data passed to Respond: %+v\n", data)

	// Ensure data is of type map[string]interface{}
	if dataMap, ok := data.(map[string]interface{}); ok {
		// Check for the template key
		if templateName, exists := dataMap["template"]; exists {
			// Ensure the template is a string
			if templateNameStr, isString := templateName.(string); isString {
				fmt.Printf("Using template: %s\n", templateNameStr) // Log the template being used
				// Remove the "template" key from the map before passing to the HTML method
				delete(dataMap, "template")
				c.HTML(status, templateNameStr, dataMap)
				return
			} else {
				fmt.Println("Error: template is not a string")
			}
		} else {
			fmt.Println("Error: template key not found")
		}
	} else {
		fmt.Println("Error: data is not a map[string]interface{}")
	}

	// Fall back to default template if no valid template is found
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
