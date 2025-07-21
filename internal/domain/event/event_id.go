package event

import (
	"fmt"

	"github.com/google/uuid"
)

type EventID struct {
	uuid.UUID
}

func NewEventID() EventID {
	id := uuid.New()
	return EventID{id}
}

func NewEventIDFromString(s string) (EventID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return EventID{}, fmt.Errorf("invalid UUID format: %w", err)
	}
	return EventID{id}, nil
}
