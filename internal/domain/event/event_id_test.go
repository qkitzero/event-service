package event

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewEventID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
	}{
		{"success new event id", true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			id := NewEventID()
			if tt.success && id.UUID == uuid.Nil {
				t.Errorf("expected valid event id, but got a nil UUID")
			}
		})
	}
}

func TestNewEventIDFromString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		id      string
	}{
		{"success new event id from string", true, "fe8c2263-bbac-4bb9-a41d-b04f5afc4425"},
		{"failure empty event id", false, ""},
		{"failure invalid event id", false, "0123456789"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewEventIDFromString(tt.id)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}
