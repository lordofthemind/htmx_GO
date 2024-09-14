package middlewares

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Process request
		c.Next()

		// Log request details
		log.Printf("Request: %s %s | Status: %d | Duration: %v",
			c.Request.Method, c.Request.URL.Path, c.Writer.Status(), time.Since(startTime))
	}
}
