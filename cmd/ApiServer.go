package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/lordofthemind/htmx_GO/internals/configs"
	"github.com/lordofthemind/htmx_GO/internals/initializers"
)

func RunApiServer() {
	log.Println("Starting server...")

	// Set up logging
	logFile, err := initializers.SetUpLoggerFile("ApiServerLogs.log")
	if err != nil {
		log.Fatalf("Failed to set up logger: %v", err)
	}
	defer logFile.Close()

	// Initialize server configuration
	err = configs.InitializeServerConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to initialize server configuration: %v", err)
	}

	// Connect to MongoDB
	ctx := context.Background()
	dsn := configs.MongoDBUrl
	timeout := 30 * time.Second
	maxRetries := 5
	dbName := "htmx_go" // Replace with your actual database name

	mongoCL, _, err := initializers.ConnectToMongoDB(ctx, dsn, timeout, maxRetries, dbName)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}
	defer mongoCL.Disconnect(ctx) // Ensure MongoDB client is disconnected on shutdown

	// Set up the Gin router with CORS
	router, err := initializers.SetUpGinServerWithCORS()
	// For without CORS use: router, err := initializers.SetUpGinServer()
	if err != nil {
		log.Fatalf("Failed to set up Gin server: %v", err)
	}

	// Create an HTTP server with the Gin router
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", configs.Port),
		Handler: router,
	}

	// Start the Gin server
	err = initializers.StartGinServer(router)
	if err != nil {
		log.Fatalf("Failed to start Gin server: %v", err)
	}

	// Handle graceful shutdown
	gracefulShutdown(server)
}

// gracefulShutdown handles the graceful shutdown of the server.
func gracefulShutdown(server *http.Server) {
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
