package models

import (
	"time"
)

// Poll represents a single poll in the system
type Poll struct {
	ID          string    `json:"id" bson:"_id"`
	Title       string    `json:"title" bson:"title"`
	Description string    `json:"description,omitempty" bson:"description,omitempty"`
	Options     []Option  `json:"options" bson:"options"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
	ExpiresAt   time.Time `json:"expires_at,omitempty" bson:"expires_at,omitempty"`
}

// Option represents a single choice in a poll
type Option struct {
	ID        string `json:"id" bson:"_id"`
	Text      string `json:"text" bson:"text"`
	VoteCount int    `json:"vote_count" bson:"vote_count"`
}

// Validate performs validation on the poll structure
func (p *Poll) Validate() error {
	// Title validation
	if p.Title == "" {
		return ErrInvalidPoll("title is required")
	}
	if len(p.Title) > 200 {
		return ErrInvalidPoll("title must be less than 200 characters")
	}

	// Description validation
	if len(p.Description) > 1000 {
		return ErrInvalidPoll("description must be less than 1000 characters")
	}

	// Options validation
	if len(p.Options) < 2 {
		return ErrInvalidPoll("poll must have at least 2 options")
	}
	if len(p.Options) > 10 {
		return ErrInvalidPoll("poll cannot have more than 10 options")
	}

	// Validate each option
	for _, option := range p.Options {
		if option.Text == "" {
			return ErrInvalidPoll("option text cannot be empty")
		}
		if len(option.Text) > 200 {
			return ErrInvalidPoll("option text must be less than 200 characters")
		}
	}

	// Expiration validation
	if !p.ExpiresAt.IsZero() && p.ExpiresAt.Before(time.Now()) {
		return ErrInvalidPoll("expiration date must be in the future")
	}

	return nil
}

// ErrInvalidPoll represents an error in poll validation
type ErrInvalidPoll string

func (e ErrInvalidPoll) Error() string {
	return string(e)
}
