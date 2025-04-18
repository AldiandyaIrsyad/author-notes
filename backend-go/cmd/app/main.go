package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	auth_service "github.com/AldiandyaIrsyad/author-notes/internal/auth"
	auth_adapter "github.com/AldiandyaIrsyad/author-notes/internal/auth/adapter"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found or could not be loaded")
	}

	mongoClient, err := connectToMongoDB()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	dbName := getEnv("MONGO_DATABASE", "author_notes")
	db := mongoClient.Database(dbName)

	authRepo := auth_adapter.NewMongoAuthRepository(db)
	authService := auth_service.NewAuthService(authRepo)
	authHandler := auth_adapter.NewAuthHTTPHandler(authService)

	router := gin.Default()

	// Use CORS middleware
	// Default allows all origins, methods, headers. Fine for development.
	// For production, configure specific origins:
	// config := cors.DefaultConfig()
	// config.AllowOrigins = []string{"http://localhost:5173", "http://your-frontend-domain.com"} // Add your frontend origin(s)
	// config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	// config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	// router.Use(cors.New(config))
	router.Use(cors.Default())

	// Register API routes
	v1 := router.Group("/v1")
	authHandler.RegisterRoutes(v1)

	// Start server
	port := getEnv("PORT", "8080")
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server starting on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func connectToMongoDB() (*mongo.Client, error) {
	uri := getEnv("MONGO_URI", "mongodb://admin:password@localhost:27017/author_notes?authSource=admin")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Ping the database to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	log.Println("Successfully connected to MongoDB")
	return client, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
