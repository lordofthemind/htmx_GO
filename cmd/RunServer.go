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
	"github.com/lordofthemind/htmx_GO/internals/tokens"
)

func RunServer() {
	log.Println("Starting server...")

	// Set up logging
	logFile, err := initializers.SetUpLoggerFile("ServerLogs.log")
	if err != nil {
		log.Fatalf("Failed to set up logger: %v", err)
	}
	defer logFile.Close()

	// Initialize server configuration
	err = configs.InitializeServerConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to initialize server configuration: %v", err)
	}

	// Set up the Gin router with optional CORS
	router, err := initializers.SetUpServerWithOptionalCORS()
	if err != nil {
		log.Fatalf("Failed to set up Gin server: %v", err)
	}

	// MongoDB connection and setup
	ctx := context.Background()
	dsn := configs.MongoDBUrl
	timeout := 30 * time.Second
	maxRetries := 5

	// Connect to MongoDB
	mongoCL, err := initializers.ConnectToMongoDB(ctx, dsn, timeout, maxRetries)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}
	defer mongoCL.Disconnect(ctx)

	// Get the MongoDB database instance
	dbName := "htmx_go"
	mongoDB := initializers.GetDatabase(mongoCL, dbName)

	// Create an HTTP server with the Gin router
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", configs.Port),
		Handler: router,
	}

	// Set up service, handler, and middleware
	repo := repositories.NewSuperuserRepository(mongoDB)
	service := services.NewSuperuserService(repo)
	tokenManager := tokens.NewJWTManager()
	handler := handlers.NewSuperuserHandler(service, tokenManager)

	// Middleware and route registration
	router.Use(middlewares.ResponseStrategyMiddleware())
	routes.RegisterSuperuserRoutes(router, handler, tokenManager)

	// Start the Gin server
	err = initializers.StartGinServer(router)
	if err != nil {
		log.Fatalf("Failed to start Gin server: %v", err)
	}

	// Handle graceful shutdown
	initializers.GracefulShutdown(server)
}
