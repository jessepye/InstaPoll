package handlers

import (
	"log"
	"net/http"
	"time"

	"instapoll/backend/models" // Import the Poll model

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo" // Import mongo driver types
	// "go.mongodb.org/mongo-driver/bson/primitive" // May be needed if using default MongoDB ObjectID
)

// PollHandler holds the database collection for polls
type PollHandler struct {
	collection *mongo.Collection // Pointer to the MongoDB collection
}

// NewPollHandler creates a new handler with the given MongoDB collection.
// This acts as a constructor for PollHandler.
func NewPollHandler(collection *mongo.Collection) *PollHandler {
	// Return a pointer to a new PollHandler instance,
	// initializing its collection field with the provided argument.
	return &PollHandler{
		collection: collection,
	}
}

// RegisterRoutes sets up the poll-related routes for the Gin engine.
// It accepts the gin.Engine directly to register the route group.
func (h *PollHandler) RegisterRoutes(r *gin.Engine) {
	// Create a route group for API endpoints prefixed with /api/polls
	polls := r.Group("/api/polls")
	{
		// Assign handler methods (which now belong to PollHandler) to specific
		// HTTP methods and paths within the group.
		polls.POST("/", h.CreatePoll)   // Handle POST requests to /api/polls
		polls.GET("/", h.ListPolls)    // Handle GET requests to /api/polls
		polls.GET("/:id", h.GetPoll)    // Handle GET requests to /api/polls/:id (with path parameter)
		// TODO: Add routes for PUT /:id (UpdatePoll) and DELETE /:id (DeletePoll) later
	}
}

// CreatePoll handles the creation of a new poll.
// It's now a method on PollHandler, allowing access to h.collection.
func (h *PollHandler) CreatePoll(c *gin.Context) {
	var poll models.Poll // Declare a variable to hold the incoming poll data

	// Attempt to bind the incoming JSON request body to the poll struct.
	if err := c.ShouldBindJSON(&poll); err != nil {
		// If binding fails (e.g., malformed JSON), return a Bad Request error.
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// --- Prepare Poll Data ---
	// Generate unique IDs for the poll and its options.
	// Note: MongoDB often uses ObjectID as the default _id. If you prefer that,
	// you wouldn't set poll.ID here and would let the driver generate it.
	// Using UUID strings as IDs is also perfectly valid.
	poll.ID = uuid.New().String() // Assign a new UUID string to the poll's ID
	for i := range poll.Options {
		poll.Options[i].ID = uuid.New().String() // Assign a new UUID string to each option's ID
		poll.Options[i].VoteCount = 0            // Initialize vote count to zero
	}
	// Set creation and update timestamps to the current time.
	now := time.Now()
	poll.CreatedAt = now
	poll.UpdatedAt = now

	// --- Validate Poll Data ---
	// Perform business logic validation using the method defined on the model.
	if err := poll.Validate(); err != nil {
		// If validation fails, return a Bad Request error with the validation message.
		log.Printf("Validation failed for poll '%s': %v", poll.Title, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed: " + err.Error()})
		return
	}

	// --- Database Interaction (TDD Step: Implement this next!) ---
	// TODO: Implement database insertion logic here.
	// This is the section that needs code to make the TestCreatePoll database assertion pass.
	/*
		Example Implementation:
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second) // Use request context with timeout
		defer cancel()

		_, err := h.collection.InsertOne(ctx, poll) // Insert the poll document
		if err != nil {
			log.Printf("Error inserting poll into database: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create poll"})
			return
		}
		log.Printf("Successfully inserted poll with ID: %s", poll.ID)
	*/

	// If the TODO block above were implemented, the poll would be saved.
	// For now, we log that it *would* be saved.
	log.Printf("TODO: Save poll to database: %+v", poll)

	// Return the newly created poll object (including generated IDs and timestamps)
	// with an HTTP 201 Created status.
	c.JSON(http.StatusCreated, poll)
}

// GetPoll handles retrieving a single poll by its ID.
// It's now a method on PollHandler.
func (h *PollHandler) GetPoll(c *gin.Context) {
	// Extract the 'id' path parameter from the URL (e.g., /api/polls/xyz).
	pollID := c.Param("id")

	// Basic validation: check if the ID parameter is empty.
	// More robust validation could check if it's a valid UUID format if needed.
	if pollID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Poll ID parameter is required"})
		return
	}

	// --- Database Interaction (TDD Step: Implement this next!) ---
	// TODO: Implement database fetch logic here using pollID.
	// This is the section that needs code to make the TestGetPoll database assertion pass.
	/*
		Example Implementation:
		var result models.Poll
		// Create a filter to find the document where the 'id' field matches pollID.
		// Note: If using MongoDB's default _id, the field name is "_id".
		// If you stored your UUID string in the 'id' field, use "id". Adjust accordingly.
		filter := bson.M{"_id": pollID} // Assuming you store your UUID in the _id field

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		// Attempt to find one document matching the filter.
		err := h.collection.FindOne(ctx, filter).Decode(&result)

		if err != nil {
			// Check if the error is because no document was found.
			if errors.Is(err, mongo.ErrNoDocuments) {
				log.Printf("Poll not found with ID: %s", pollID)
				c.JSON(http.StatusNotFound, gin.H{"error": "Poll not found"})
			} else {
				// Handle other potential database errors.
				log.Printf("Error retrieving poll with ID %s: %v", pollID, err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve poll"})
			}
			return
		}

		// If found, return the poll data with HTTP 200 OK.
		c.JSON(http.StatusOK, result)
		return // Important to return here after successful handling
	*/

	// If the TODO block above were implemented, the poll would be fetched.
	// For now, we log that it *would* be fetched.
	log.Printf("TODO: Fetch poll with ID %s from database", pollID)

	// Return a 404 Not Found status as the database logic isn't implemented yet.
	c.JSON(http.StatusNotFound, gin.H{"error": "Poll not found (DB fetch not implemented)"})
}

// ListPolls handles retrieving a list of all polls.
// It's now a method on PollHandler. Pagination should be added later.
func (h *PollHandler) ListPolls(c *gin.Context) {
	// --- Database Interaction (TDD Step: Implement this next!) ---
	// TODO: Implement database fetch logic here to retrieve multiple documents.
	// This is the section that needs code to make the TestListPolls database assertion pass.
	/*
		Example Implementation:
		var results []models.Poll // Slice to hold the results

		// Create an empty filter {} to match all documents.
		filter := bson.M{}

		// Add query parameters for pagination later (e.g., limit, skip/offset)
		// findOptions := options.Find()
		// findOptions.SetLimit(10) // Example limit

		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second) // Longer timeout for potentially larger lists
		defer cancel()

		// Find documents matching the filter.
		cursor, err := h.collection.Find(ctx, filter) // Add findOptions here if using pagination
		if err != nil {
			log.Printf("Error finding polls: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve polls"})
			return
		}
		// Ensure the cursor is closed when the function returns.
		defer cursor.Close(ctx)

		// Decode all documents found by the cursor into the results slice.
		if err = cursor.All(ctx, &results); err != nil {
			log.Printf("Error decoding polls from cursor: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode polls"})
			return
		}

		// Important: If no documents are found, cursor.All returns an empty slice and no error.
		// However, if the results slice was initially nil, encoding it to JSON might result in 'null'.
		// Ensure we return an empty JSON array '[]' instead of 'null' if no polls exist.
		if results == nil {
			results = []models.Poll{}
		}

		// Return the list of polls with HTTP 200 OK.
		c.JSON(http.StatusOK, results)
		return // Important to return here after successful handling
	*/

	// If the TODO block above were implemented, the polls would be listed.
	// For now, we log that it *would* happen.
	log.Println("TODO: Fetch list of polls from database")

	// Return an empty list with HTTP 200 OK as the DB logic isn't implemented yet.
	c.JSON(http.StatusOK, []models.Poll{})
}
