package event

import (
	"time"

	"github.com/qkitzero/event-service/internal/domain/event"
)

type EventModel struct {
	ID          event.EventID
	Title       event.Title
	Description event.Description
	CreatedAt   time.Time
}

func (EventModel) TableName() string {
	return "events"
}
