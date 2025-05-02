package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"instapoll/backend/models" // Import models

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require" // Use require for fatal assertions in setup
	"go.mongodb.org/mongo-driver/bson"    // Import bson
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	testMongoURI    = "mongodb://localhost:27017" // Default for local Docker
	testDatabase    = "instapoll_test"            // Separate DB for tests
	testCollection  = "polls"
	defaultTimeout  = 10 * time.Second
)

var testPollCollection *mongo.Collection // Make collection accessible to tests

// setupTestDB connects to the test database and returns the client and collection
func setupTestDB(t *testing.T) (*mongo.Client, *mongo.Collection) {
	// Use require for setup steps, as failure here means tests can't run
	mongoURI := os.Getenv("TEST_MONGODB_URI") // Allow overriding via env var
	if mongoURI == "" {
		mongoURI = testMongoURI
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	require.NoError(t, err, "Failed to connect to test MongoDB")

	err = client.Ping(ctx, readpref.Primary())
	require.NoError(t, err, "Failed to ping test MongoDB")

	collection := client.Database(testDatabase).Collection(testCollection)
	testPollCollection = collection // Assign to global var for test access
	log.Println("Connected to test database:", testDatabase)
	return client, collection
}

// teardownTestDB disconnects from the database and cleans up
func teardownTestDB(t *testing.T, client *mongo.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	// Clean up the test database
	err := client.Database(testDatabase).Drop(ctx)
	// Use assert here instead of require, so teardown continues even if drop fails
	assert.NoError(t, err, "Failed to drop test database")
	log.Println("Dropped test database:", testDatabase)


	err = client.Disconnect(ctx)
	assert.NoError(t, err, "Failed to disconnect from test MongoDB")
	log.Println("Disconnected from test database.")
}

// setupRouter creates a Gin router with the test handler and DB collection
func setupRouter(collection *mongo.Collection) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	// *** This line (84) causes 'undefined: NewPollHandler' if poll.go is incorrect ***
	pollHandler := NewPollHandler(collection) // Create handler with test collection
	pollHandler.RegisterRoutes(r)             // Register routes
	return r
}

// TestMain runs setup and teardown for the entire test suite
func TestMain(m *testing.M) {
	// Setup: Connect to DB
	// *** FIX: Use blank identifier '_' for the collection returned here ***
	// *** as the local variable 'collection' was unused in this scope.  ***
	// *** The global 'testPollCollection' is assigned inside setupTestDB. ***
	client, _ := setupTestDB(&testing.T{}) // Use dummy T for setup

	// Run tests
	exitCode := m.Run()

	// Teardown: Disconnect and clean DB
	teardownTestDB(&testing.T{}, client) // Use dummy T for teardown

	os.Exit(exitCode)
}

// clearTestCollection removes all documents from the test collection
func clearTestCollection(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	_, err := testPollCollection.DeleteMany(ctx, bson.M{})
	require.NoError(t, err, "Failed to clear test collection")
}


// --- Test Cases ---

func TestCreatePoll(t *testing.T) {
	// Ensure collection is clean before test
	clearTestCollection(t)

	// Setup router with the globally available testPollCollection
	router := setupRouter(testPollCollection)

	validPollPayload := models.Poll{
		Title:       "Favorite Color?",
		Description: "A simple poll",
		Options: []models.Option{
			{Text: "Red"},
			{Text: "Blue"},
		},
		// ID, Timestamps, Option IDs/Votes are set by the handler
	}

	jsonData, err := json.Marshal(validPollPayload)
	require.NoError(t, err) // Use require for setup steps within the test

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/polls", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	// Assert HTTP Status
	assert.Equal(t, http.StatusCreated, w.Code, "Expected status code 201 Created")

	// Assert Response Body (check if basic structure is returned)
	var responsePoll models.Poll
	err = json.Unmarshal(w.Body.Bytes(), &responsePoll)
	require.NoError(t, err, "Failed to unmarshal response body")
	assert.NotEmpty(t, responsePoll.ID, "Response poll ID should not be empty")
	assert.Equal(t, validPollPayload.Title, responsePoll.Title, "Response title mismatch")
	assert.Len(t, responsePoll.Options, 2, "Response should have 2 options")
	assert.NotEmpty(t, responsePoll.Options[0].ID, "Response option 0 ID should not be empty")
	assert.NotEmpty(t, responsePoll.Options[1].ID, "Response option 1 ID should not be empty")
	assert.NotZero(t, responsePoll.CreatedAt, "Response CreatedAt should be set")
	assert.NotZero(t, responsePoll.UpdatedAt, "Response UpdatedAt should be set")

	// --- TDD Failing Assertion ---
	// Verify the poll was actually saved to the database
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	var dbPoll models.Poll
	// Use the ID returned in the response to find the document in the DB
	filter := bson.M{"_id": responsePoll.ID} // Assuming ID is used as _id
	err = testPollCollection.FindOne(ctx, filter).Decode(&dbPoll)

	// This assertion WILL FAIL until InsertOne is implemented in the handler
	assert.NoError(t, err, "Poll should be found in the database")
	if err == nil { // Only compare if found
		assert.Equal(t, responsePoll.Title, dbPoll.Title, "Database poll title mismatch")
		assert.Equal(t, len(responsePoll.Options), len(dbPoll.Options), "Database poll options count mismatch")
	}
}

func TestGetPoll(t *testing.T) {
	// Ensure collection is clean before test
	clearTestCollection(t)

	// Setup router with the globally available testPollCollection
	router := setupRouter(testPollCollection)

	// --- TDD Setup: Insert a poll directly into the DB ---
	pollID := uuid.New().String()
	testPoll := models.Poll{
		ID:          pollID,
		Title:       "Existing Poll",
		Description: "This poll exists for testing GET",
		Options: []models.Option{
			{ID: uuid.New().String(), Text: "Yes", VoteCount: 5},
			{ID: uuid.New().String(), Text: "No", VoteCount: 2},
		},
		CreatedAt: time.Now().Add(-time.Hour), // Set explicit times
		UpdatedAt: time.Now().Add(-time.Minute),
	}
	ctxInsert, cancelInsert := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancelInsert()
	_, err := testPollCollection.InsertOne(ctxInsert, testPoll)
	require.NoError(t, err, "Failed to insert test poll directly into DB")

	// --- Make API Request ---
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/polls/"+pollID, nil)
	router.ServeHTTP(w, req)

	// --- TDD Failing Assertions ---
	// This assertion WILL FAIL until FindOne is implemented in the handler
	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200 OK")

	if w.Code == http.StatusOK {
		var responsePoll models.Poll
		err = json.Unmarshal(w.Body.Bytes(), &responsePoll)
		require.NoError(t, err, "Failed to unmarshal response body")
		assert.Equal(t, testPoll.ID, responsePoll.ID, "Response ID mismatch")
		assert.Equal(t, testPoll.Title, responsePoll.Title, "Response title mismatch")
		assert.Len(t, responsePoll.Options, 2, "Response options count mismatch")
		// Optionally compare timestamps more precisely if needed
		// assert.WithinDuration(t, testPoll.CreatedAt, responsePoll.CreatedAt, time.Second)
	}
}

func TestGetPoll_NotFound(t *testing.T) {
	// Ensure collection is clean
	clearTestCollection(t)
	router := setupRouter(testPollCollection)

	nonExistentID := uuid.New().String()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/polls/"+nonExistentID, nil)
	router.ServeHTTP(w, req)

	// This assertion SHOULD PASS even without DB logic IF the handler correctly
	// returns 404 when the (future) DB lookup fails.
	// However, the current stub returns 404 always, so it passes for the wrong reason initially.
	assert.Equal(t, http.StatusNotFound, w.Code, "Expected status code 404 Not Found")
}


func TestListPolls(t *testing.T) {
	// Ensure collection is clean
	clearTestCollection(t)
	router := setupRouter(testPollCollection)

	// --- TDD Setup: Insert multiple polls directly into the DB ---
	poll1 := models.Poll{ID: uuid.New().String(), Title: "Poll 1", Options: []models.Option{{Text: "A"}, {Text: "B"}}, CreatedAt: time.Now()}
	poll2 := models.Poll{ID: uuid.New().String(), Title: "Poll 2", Options: []models.Option{{Text: "C"}, {Text: "D"}}, CreatedAt: time.Now()}

	ctxInsert, cancelInsert := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancelInsert()
	_, err := testPollCollection.InsertMany(ctxInsert, []interface{}{poll1, poll2})
	require.NoError(t, err, "Failed to insert multiple test polls")

	// --- Make API Request ---
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/polls", nil)
	router.ServeHTTP(w, req)

	// --- TDD Failing Assertions ---
	// These assertions WILL FAIL until Find is implemented in the handler
	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200 OK")

	if w.Code == http.StatusOK {
		var responsePolls []models.Poll
		err = json.Unmarshal(w.Body.Bytes(), &responsePolls)
		require.NoError(t, err, "Failed to unmarshal response body")
		assert.Len(t, responsePolls, 2, "Expected 2 polls in the response list")
		// Optionally, check if the returned polls match the inserted ones
		// This might require sorting or checking IDs if order isn't guaranteed
	}
}

func TestListPolls_Empty(t *testing.T) {
	// Ensure collection is clean
	clearTestCollection(t)
	router := setupRouter(testPollCollection)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/polls", nil)
	router.ServeHTTP(w, req)

	// This assertion SHOULD PASS even without DB logic IF the handler correctly
	// returns an empty array `[]` when the (future) DB lookup finds nothing.
	// The current stub returns `[]`, so it passes for the wrong reason initially.
	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200 OK")
	assert.Equal(t, "[]", w.Body.String(), "Expected empty JSON array '[]'")
}


// Add tests for invalid CreatePoll payloads (these should still work as they test validation before DB)
func TestCreatePoll_InvalidPayload(t *testing.T) {
	clearTestCollection(t)
	router := setupRouter(testPollCollection)

	tests := []struct {
		name       string
		payload    string // Send raw JSON string to test binding/validation errors
		wantStatus int
	}{
		{
			name:       "invalid json",
			payload:    `{"title": "Test", "options": [`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "empty title",
			payload:    `{"title": "", "options": [{"text":"A"},{"text":"B"}]}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "single option",
			payload:    `{"title": "Test", "options": [{"text":"A"}]}`,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/polls", bytes.NewBufferString(tt.payload))
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(w, req)
			assert.Equal(t, tt.wantStatus, w.Code)

			// Verify DB was NOT touched for bad requests
			ctxCount, cancelCount := context.WithTimeout(context.Background(), defaultTimeout)
			defer cancelCount()
			count, err := testPollCollection.CountDocuments(ctxCount, bson.M{})
			require.NoError(t, err)
			assert.Zero(t, count, "Database should be empty after invalid request")
		})
	}
}

