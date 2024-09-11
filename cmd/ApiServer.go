package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/lordofthemind/htmx_GO/internals/configs"
	"github.com/lordofthemind/htmx_GO/internals/initializers"
)

func RunApiServer() {
	log.Println("Starting server...")

	logFile, err := initializers.SetUpLoggerFile("ApiServerLogs.log")
	if err != nil {
		log.Fatalf("Failed to set up logger: %v", err)
	}
	defer logFile.Close()

	err = configs.InitializeServerConfig("sGin.yaml")
	if err != nil {
		log.Fatalf("Failed to set up logger: %v", err)
	}

	ctx := context.Background()
	dsn := configs.MongoDBUrl
	timeout := 30 * time.Second
	maxRetries := 5
	dbName := "htmx_go" // Replace with your actual database name

	mongoCL, mongoDB, err := initializers.ConnectToMongoDB(ctx, dsn, timeout, maxRetries, dbName)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	// Create a Gin router (with or without CORS as needed)
	router, err := initializers.SetUpGinServerWithCORS()
	// router, err := initializers.SetUpGinServer()
	if err != nil {
		log.Fatalf("Failed to set up Gin server: %v", err)
	}

	// Start the Gin server
	if err := initializers.StartGinServer(router); err != nil {
		log.Fatalf("Failed to start Gin server: %v", err)
	}

	// Create an HTTP server to handle graceful shutdown
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", configs.Port),
		Handler: router,
	}

	// Graceful shutdown handling
	if err := initializers.StopGinServer(server); err != nil {
		log.Fatalf("Failed to stop Gin server: %v", err)
	}
}
