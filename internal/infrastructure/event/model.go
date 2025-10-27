package event

import (
	"time"

	"github.com/qkitzero/event-service/internal/domain/event"
	"github.com/qkitzero/event-service/internal/domain/user"
)

type EventModel struct {
	ID          event.EventID
	UserID      user.UserID
	Title       event.Title
	Description event.Description
	StartTime   time.Time
	EndTime     time.Time
	Color       event.Color
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (EventModel) TableName() string {
	return "events"
}
