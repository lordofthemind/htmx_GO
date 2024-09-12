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

func (r *HTMLResponseStrategy) Respond(c *gin.Context, data interface{}, status int) {
	fmt.Printf("Data passed to Respond: %+v\n", data)

	dataMap, ok := data.(map[string]interface{})
	if !ok {
		fmt.Println("Error: data is not a map[string]interface{}")
		c.HTML(status, "error.html", gin.H{"error": "Internal Server Error"})
		return
	}

	templateName, exists := dataMap["template"]
	if !exists {
		fmt.Println("Error: template key not found in data")
		c.HTML(status, "error.html", gin.H{"error": "Template not specified"})
		return
	}

	templateNameStr, isString := templateName.(string)
	if !isString {
		fmt.Println("Error: template name is not a string")
		c.HTML(status, "error.html", gin.H{"error": "Invalid template name"})
		return
	}

	fmt.Printf("Using template: %s\n", templateNameStr)
	delete(dataMap, "template")
	c.HTML(status, templateNameStr, dataMap)
}

// func (r *HTMLResponseStrategy) Respond(c *gin.Context, data interface{}, status int) {
// 	// Log the data for debugging
// 	fmt.Printf("Data passed to Respond: %+v\n", data)

// 	// Check if data is a map[string]interface{}
// 	dataMap, ok := data.(map[string]interface{})
// 	if !ok {
// 		fmt.Println("Error: data is not a map[string]interface{}")
// 		c.HTML(status, "error.html", gin.H{"error": "Internal Server Error"})
// 		return
// 	}

// 	// Check if template key exists and is a string
// 	templateName, exists := dataMap["template"]
// 	if !exists {
// 		fmt.Println("Error: template key not found in data")
// 		c.HTML(status, "error.html", gin.H{"error": "Template not specified"})
// 		return
// 	}

// 	templateNameStr, isString := templateName.(string)
// 	if !isString {
// 		fmt.Println("Error: template name is not a string")
// 		c.HTML(status, "error.html", gin.H{"error": "Invalid template name"})
// 		return
// 	}

// 	fmt.Printf("Using template: %s\n", templateNameStr)
// 	// Remove the "template" key before passing data to c.HTML
// 	delete(dataMap, "template")
// 	c.HTML(status, templateNameStr, dataMap)
// }

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
