package handlers

import (
	"net/http"
	"time"

	"instapoll/backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreatePoll handles the creation of a new poll
func CreatePoll(c *gin.Context) {
	var poll models.Poll
	if err := c.ShouldBindJSON(&poll); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate IDs and set timestamps
	poll.ID = uuid.New().String()
	for i := range poll.Options {
		poll.Options[i].ID = uuid.New().String()
		poll.Options[i].VoteCount = 0
	}
	poll.CreatedAt = time.Now()
	poll.UpdatedAt = time.Now()

	// Validate the poll
	if err := poll.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Save to database
	// For now, just return the created poll
	c.JSON(http.StatusCreated, poll)
}

// GetPoll handles retrieving a single poll by ID
func GetPoll(c *gin.Context) {
	pollID := c.Param("id")
	if pollID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "poll ID is required"})
		return
	}

	// TODO: Fetch from database
	// For now, return a 404
	c.JSON(http.StatusNotFound, gin.H{"error": "poll not found"})
}

// ListPolls handles retrieving all polls
func ListPolls(c *gin.Context) {
	// TODO: Fetch from database with pagination
	// For now, return empty list
	c.JSON(http.StatusOK, []models.Poll{})
}

// RegisterPollRoutes sets up the poll-related routes
func RegisterPollRoutes(r *gin.Engine) {
	polls := r.Group("/api/polls")
	{
		polls.POST("/", CreatePoll)
		polls.GET("/", ListPolls)
		polls.GET("/:id", GetPoll)
	}
}
