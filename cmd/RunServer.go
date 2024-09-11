package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/lordofthemind/htmx_GO/internals/configs"
	"github.com/lordofthemind/htmx_GO/internals/handlers"
	"github.com/lordofthemind/htmx_GO/internals/initializers"
	"github.com/lordofthemind/htmx_GO/internals/middlewares"
	"github.com/lordofthemind/htmx_GO/internals/repositories"
	"github.com/lordofthemind/htmx_GO/internals/routes"
	"github.com/lordofthemind/htmx_GO/internals/services"
)

func RunServer() {
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

	mongoCL, mongoDB, err := initializers.ConnectToMongoDB(ctx, dsn, timeout, maxRetries, dbName)
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

	// In runServer.go
	repo := repositories.NewSuperuserRepository(mongoDB)
	service := services.NewSuperuserService(repo)
	handler := handlers.NewSuperuserHandler(service)

	// In your `RunServer()` function
	router.Use(middlewares.ResponseStrategyMiddleware())

	// Register routes
	routes.RegisterSuperuserRoutes(router, handler)

	// Handle graceful shutdown
	initializers.GracefulShutdown(server)
}
