package main

import (
	"context" // Required for database operations
	"log"     // For logging messages
	"net/http"
	"os"      // To read environment variables
	"time"    // For setting timeouts

	// Import the handlers package from the current module
	"instapoll/backend/handlers"

	"github.com/gin-gonic/gin"                // Gin web framework
	"go.mongodb.org/mongo-driver/mongo"       // MongoDB Go Driver
	"go.mongodb.org/mongo-driver/mongo/options" // MongoDB Driver options
	"go.mongodb.org/mongo-driver/mongo/readpref"  // For pinging the database
)

// Constants for database configuration
const (
	// Default MongoDB connection string (used if MONGODB_URI env var is not set)
	defaultMongoURI = "mongodb://localhost:27017"
	// Name of the database to use within MongoDB
	databaseName = "instapoll"
	// Name of the collection to store poll documents
	collectionName = "polls"
	// Timeout duration for database operations like connect/ping
	dbTimeout = 10 * time.Second
)

func main() {
	log.Println("Starting InstaPoll backend service...")

	// --- Database Connection Setup ---
	log.Println("Attempting to connect to MongoDB...")

	// Get MongoDB connection URI from environment variable MONGODB_URI.
	// Fallback to the default URI if the environment variable is not set.
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = defaultMongoURI
		log.Printf("MONGODB_URI environment variable not set, using default: %s", mongoURI)
	}

	// Create a context with a timeout for the database connection attempt.
	// This prevents the application from hanging indefinitely if the DB is unavailable.
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	// Ensure the context resources are released when main() exits.
	defer cancel()

	// Configure the MongoDB client options using the URI.
	clientOptions := options.Client().ApplyURI(mongoURI)
	// Connect to the MongoDB server.
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		// If connection fails, log the error and exit fatally.
		log.Fatalf("FATAL: Failed to connect to MongoDB at %s: %v", mongoURI, err)
	}

	// Ping the primary node of the MongoDB cluster to verify the connection is active.
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		// If ping fails, log the error and exit fatally.
		log.Fatalf("FATAL: Failed to ping MongoDB: %v", err)
	}
	log.Println("Successfully connected and pinged MongoDB.")

	// Set up a deferred function to disconnect from MongoDB when the main function exits.
	// This ensures graceful shutdown.
	defer func() {
		log.Println("Disconnecting from MongoDB...")
		// Use a background context for disconnection as the original context might have expired.
		disconnectCtx, disconnectCancel := context.WithTimeout(context.Background(), dbTimeout)
		defer disconnectCancel()
		if err = client.Disconnect(disconnectCtx); err != nil {
			// Log any errors during disconnection.
			log.Printf("Error disconnecting from MongoDB: %v", err)
		} else {
			log.Println("Successfully disconnected from MongoDB.")
		}
	}()

	// Get a handle for the specific database ("instapoll").
	db := client.Database(databaseName)
	// Get a handle for the specific collection ("polls") within the database.
	pollCollection := db.Collection(collectionName)
	log.Printf("Using database '%s' and collection '%s'", databaseName, collectionName)

	// --- Gin Router and Handler Setup ---
	log.Println("Setting up Gin router and routes...")
	// Create a new Gin engine with default middleware (logger, recovery).
	r := gin.Default()

	// Create an instance of PollHandler, passing the database collection handle.
	// This injects the database dependency into the handler.
	pollHandler := handlers.NewPollHandler(pollCollection)

	// Register the API routes defined in the PollHandler.
	// This calls the RegisterRoutes method on the pollHandler instance.
	pollHandler.RegisterRoutes(r)
	log.Println("Registered poll routes under /api/polls")

	// Define a simple root endpoint for health checks or basic info.
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "InstaPoll Backend is running!",
			"status":  "OK",
		})
	})

	// --- Start HTTP Server ---
	serverAddr := ":8080" // Address and port to listen on
	log.Printf("Starting HTTP server, listening on %s", serverAddr)
	// Start the Gin server and listen for incoming requests.
	// r.Run() blocks until the server is shut down or an error occurs.
	if err := r.Run(serverAddr); err != nil && err != http.ErrServerClosed {
		// Log fatal error if the server fails to start (excluding graceful shutdown).
		log.Fatalf("FATAL: Failed to run server: %v", err)
	}

	log.Println("Server shut down gracefully.")
}
