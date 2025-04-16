package models

import (
	"testing"
	"time"
)

func TestPollValidation(t *testing.T) {
	tests := []struct {
		name    string
		poll    Poll
		wantErr bool
	}{
		{
			name: "valid poll",
			poll: Poll{
				Title: "Test Poll",
				Options: []Option{
					{Text: "Option 1"},
					{Text: "Option 2"},
				},
			},
			wantErr: false,
		},
		{
			name: "empty title",
			poll: Poll{
				Title: "",
				Options: []Option{
					{Text: "Option 1"},
					{Text: "Option 2"},
				},
			},
			wantErr: true,
		},
		{
			name: "title too long",
			poll: Poll{
				Title: "This is a very long title that exceeds the maximum allowed length of 200 characters. This is a very long title that exceeds the maximum allowed length of 200 characters. This is a very long title that exceeds the maximum allowed length of 200 characters.",
				Options: []Option{
					{Text: "Option 1"},
					{Text: "Option 2"},
				},
			},
			wantErr: true,
		},
		{
			name: "single option",
			poll: Poll{
				Title: "Test Poll",
				Options: []Option{
					{Text: "Option 1"},
				},
			},
			wantErr: true,
		},
		{
			name: "too many options",
			poll: Poll{
				Title: "Test Poll",
				Options: []Option{
					{Text: "Option 1"}, {Text: "Option 2"}, {Text: "Option 3"},
					{Text: "Option 4"}, {Text: "Option 5"}, {Text: "Option 6"},
					{Text: "Option 7"}, {Text: "Option 8"}, {Text: "Option 9"},
					{Text: "Option 10"}, {Text: "Option 11"},
				},
			},
			wantErr: true,
		},
		{
			name: "empty option text",
			poll: Poll{
				Title: "Test Poll",
				Options: []Option{
					{Text: ""},
					{Text: "Option 2"},
				},
			},
			wantErr: true,
		},
		{
			name: "expired poll",
			poll: Poll{
				Title: "Test Poll",
				Options: []Option{
					{Text: "Option 1"},
					{Text: "Option 2"},
				},
				ExpiresAt: time.Now().Add(-24 * time.Hour),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.poll.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Poll.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
