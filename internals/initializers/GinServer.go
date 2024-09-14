package initializers

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lordofthemind/htmx_GO/internals/configs"
	"github.com/lordofthemind/htmx_GO/pkgs/helpers"
)

// ServerSetup defines the interface for setting up a Gin server
type ServerSetup interface {
	SetUpServer() (*gin.Engine, error)
}

// CorsServerSetup is the struct for setting up the server with CORS
type CorsServerSetup struct{}

// SetUpServer sets up a Gin server with CORS middleware
func (c *CorsServerSetup) SetUpServer() (*gin.Engine, error) {
	router := gin.Default()

	// Serve static files using the path from the config
	router.Static("/static", configs.StaticPath)

	// Load HTML templates using the path from the config
	router.LoadHTMLGlob(configs.TemplatePath)

	// Configure CORS using application-specific settings
	config := cors.Config{
		AllowOrigins:     configs.AllowedOrigins,
		AllowMethods:     configs.AllowedMethods,
		AllowHeaders:     configs.AllowedHeaders,
		ExposeHeaders:    configs.ExposedHeaders,
		AllowCredentials: configs.AllowedCredentials,
	}

	// Apply CORS middleware
	router.Use(cors.New(config))

	// Log CORS settings for debugging
	log.Printf("CORS configured with origins: %v, methods: %v, headers: %v, expose headers: %v, allow credentials: %v",
		config.AllowOrigins, config.AllowMethods, config.AllowHeaders, config.ExposeHeaders, config.AllowCredentials)

	return router, nil
}

// BasicServerSetup is the struct for setting up the server without CORS
type BasicServerSetup struct{}

// SetUpServer sets up a basic Gin server without CORS
func (b *BasicServerSetup) SetUpServer() (*gin.Engine, error) {
	router := gin.Default()

	// Serve static files using the path from the config
	router.Static("/static", configs.StaticPath)

	// Load HTML templates using the path from the config
	router.LoadHTMLGlob(configs.TemplatePath)

	return router, nil
}

// SetUpServerWithOptionalCORS sets up the Gin router with or without CORS based on the UseCORS flag.
func SetUpServerWithOptionalCORS() (*gin.Engine, error) {
	var serverSetup ServerSetup

	// Choose the server setup based on UseCORS config
	if configs.UseCORS {
		fmt.Println("Setting up server with CORS...")
		serverSetup = &CorsServerSetup{}
	} else {
		fmt.Println("Setting up server without CORS...")
		serverSetup = &BasicServerSetup{}
	}

	// Set up the Gin router
	return serverSetup.SetUpServer()
}

// StartGinServer starts the provided Gin server with or without TLS based on application settings.
// It returns an error if the server fails to start or if there are issues loading the TLS certificates.
func StartGinServer(router *gin.Engine) error {
	// Create the HTTP server configuration with Gin as the handler
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", configs.Port),
		Handler: router,
	}

	// Check if TLS is enabled and configure the server accordingly
	if configs.UseTLS {
		// Load the TLS certificate and key
		cert, err := helpers.LoadTLSCertificate(configs.TlsCertFile, configs.TlsKeyFile)
		if err != nil {
			return fmt.Errorf("failed to load TLS certificate: %w", err)
		}

		// Configure the server with TLS settings
		server.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
		log.Printf("Gin server is running on port %d with TLS", configs.Port)

		// Start the server with TLS
		go func() {
			if err := server.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
				log.Printf("ListenAndServeTLS error: %v", err)
			}
		}()
	} else {
		log.Printf("Gin server is running on port %d without TLS", configs.Port)

		// Start the server without TLS
		go func() {
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Printf("ListenAndServe error: %v", err)
			}
		}()
	}

	return nil
}

// GracefulShutdown handles the graceful shutdown of the server.
func GracefulShutdown(server *http.Server) {
	// Create a channel to listen for interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down server...")

	// Context with timeout for shutdown
	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to shut down the server gracefully
	if err := server.Shutdown(ctxShutDown); err != nil {
		log.Printf("Server forced to shutdown: %v", err)

		// Retry mechanism for shutdown
		for retries := 0; retries < 3; retries++ {
			log.Printf("Retrying shutdown... attempt %d", retries+1)
			if err := server.Shutdown(ctxShutDown); err == nil {
				log.Println("Server shutdown successfully on retry")
				return
			}
		}
		log.Fatalf("Failed to shutdown server gracefully after retries: %v", err)
	}

	log.Println("Server shutdown successfully")
}
