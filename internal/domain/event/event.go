package event

import (
	"time"

	"github.com/qkitzero/event-service/internal/domain/user"
)

type Event interface {
	ID() EventID
	UserID() user.UserID
	Title() Title
	Description() Description
	StartTime() time.Time
	EndTime() time.Time
	CreatedAt() time.Time
}

type event struct {
	id          EventID
	userID      user.UserID
	title       Title
	description Description
	startTime   time.Time
	endTime     time.Time
	createdAt   time.Time
}

func (e event) ID() EventID {
	return e.id
}

func (e event) UserID() user.UserID {
	return e.userID
}

func (e event) Title() Title {
	return e.title
}

func (e event) Description() Description {
	return e.description
}

func (e event) StartTime() time.Time {
	return e.startTime
}

func (e event) EndTime() time.Time {
	return e.endTime
}

func (e event) CreatedAt() time.Time {
	return e.createdAt
}

func NewEvent(
	id EventID,
	userID user.UserID,
	title Title,
	description Description,
	startTime time.Time,
	endTime time.Time,
	createdAt time.Time,
) Event {
	return &event{
		id:          id,
		userID:      userID,
		title:       title,
		description: description,
		startTime:   startTime,
		endTime:     endTime,
		createdAt:   createdAt,
	}
}
