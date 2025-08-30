package event

import (
	"testing"
	"time"

	"github.com/qkitzero/event-service/internal/domain/user"
)

func TestNewEvent(t *testing.T) {
	t.Parallel()
	id, err := NewEventIDFromString("fe8c2263-bbac-4bb9-a41d-b04f5afc4425")
	if err != nil {
		t.Errorf("failed to new event id: %v", err)
	}
	userID, err := user.NewUserIDFromString("6d322c66-bf4d-427a-970c-874f3745f653")
	if err != nil {
		t.Errorf("failed to new user id: %v", err)
	}
	title, err := NewTitle("title")
	if err != nil {
		t.Errorf("failed to new title: %v", err)
	}
	description, err := NewDescription("description")
	if err != nil {
		t.Errorf("failed to new description: %v", err)
	}
	tests := []struct {
		name        string
		success     bool
		id          EventID
		userID      user.UserID
		title       Title
		description Description
		startTime   time.Time
		endTime     time.Time
		createdAt   time.Time
	}{
		{"success new event", true, id, userID, title, description, time.Now(), time.Now(), time.Now()},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			event := NewEvent(tt.id, tt.userID, tt.title, tt.description, tt.startTime, tt.endTime, tt.createdAt)
			if tt.success && event.ID() != tt.id {
				t.Errorf("ID() = %v, want %v", event.ID(), tt.id)
			}
			if tt.success && event.UserID() != tt.userID {
				t.Errorf("UserID() = %v, want %v", event.UserID(), tt.userID)
			}
			if tt.success && event.Title() != tt.title {
				t.Errorf("Title() = %v, want %v", event.Title(), tt.title)
			}
			if tt.success && event.Description() != tt.description {
				t.Errorf("Description() = %v, want %v", event.Description(), tt.description)
			}
			if tt.success && !event.StartTime().Equal(tt.startTime) {
				t.Errorf("StartTime() = %v, want %v", event.StartTime(), tt.startTime)
			}
			if tt.success && !event.EndTime().Equal(tt.endTime) {
				t.Errorf("EndTime() = %v, want %v", event.EndTime(), tt.endTime)
			}
			if tt.success && !event.CreatedAt().Equal(tt.createdAt) {
				t.Errorf("CreatedAt() = %v, want %v", event.CreatedAt(), tt.createdAt)
			}
		})
	}
}
