package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"instapoll/backend/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	RegisterPollRoutes(r)
	return r
}

func TestCreatePoll(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name       string
		payload    models.Poll
		wantStatus int
	}{
		{
			name: "valid poll",
			payload: models.Poll{
				Title: "Test Poll",
				Options: []models.Option{
					{Text: "Option 1"},
					{Text: "Option 2"},
				},
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "invalid poll - empty title",
			payload: models.Poll{
				Title: "",
				Options: []models.Option{
					{Text: "Option 1"},
					{Text: "Option 2"},
				},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid poll - single option",
			payload: models.Poll{
				Title: "Test Poll",
				Options: []models.Option{
					{Text: "Option 1"},
				},
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.payload)
			assert.NoError(t, err)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/polls", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(w, req)
			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantStatus == http.StatusCreated {
				var response models.Poll
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotEmpty(t, response.ID)
				assert.Equal(t, tt.payload.Title, response.Title)
				assert.Len(t, response.Options, len(tt.payload.Options))
				for i, option := range response.Options {
					assert.NotEmpty(t, option.ID)
					assert.Equal(t, tt.payload.Options[i].Text, option.Text)
					assert.Equal(t, 0, option.VoteCount)
				}
			}
		})
	}
}

func TestGetPoll(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name       string
		pollID     string
		wantStatus int
	}{
		{
			name:       "non-existent poll",
			pollID:     "non-existent-id",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "empty poll ID",
			pollID:     "",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/polls/"+tt.pollID, nil)

			router.ServeHTTP(w, req)
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestListPolls(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/polls", nil)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.Poll
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Empty(t, response) // Currently returns empty list
}
