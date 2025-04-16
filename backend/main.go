package main

import (
	"net/http"

	"instapoll/backend/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Create a new Gin router
	r := gin.Default()

	// Register poll routes
	handlers.RegisterPollRoutes(r)

	// Define a simple "hello world" endpoint
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello from InstaPoll!",
		})
	})

	// Start the server on port 8080
	r.Run(":8080")
}
