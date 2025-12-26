package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-auth-service/src/app"
	"github.com/go-auth-service/src/config"
	"github.com/go-auth-service/src/domain"
	"github.com/go-auth-service/src/routes"
	"github.com/joho/godotenv"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env vars")
	}

	// Connect to MongoDB
	client, err := config.ConnectDB(ctx)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	DB := client.Database("auth")
	if DB == nil {
		log.Fatal("Mongo database is nil after connection")
	}
	domain.UsersCollection = DB.Collection("users")
	log.Println("MongoDB connected...")

	// Initialize Fiber server
	server := app.NewServer()

	// Setup routes
	api := server.Group("/api")
	routes.AuthHandler(api)

	// Start server in a goroutine
	go func() {
		if err := server.Listen(":" + config.Load().Port); err != nil {
			log.Println("Server stopped:", err)
		}
	}()

	log.Println("Server started on port", config.Load().Port)

	// Wait for Ctrl+C or SIGTERM
	<-ctx.Done()
	log.Println("Shutting down server...")

	// Give Fiber 5 seconds to shutdown gracefully
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(); err != nil {
		log.Println("Error during server shutdown:", err)
	}

	// Disconnect MongoDB
	if err := client.Disconnect(shutdownCtx); err != nil {
		log.Println("Error disconnecting MongoDB:", err)
	}

	log.Println("Server exited properly")
}
