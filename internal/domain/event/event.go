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
	Color() Color
	CreatedAt() time.Time
	UpdatedAt() time.Time
	Update(title Title, description Description, startTime, endTime time.Time, color Color)
}

type event struct {
	id          EventID
	userID      user.UserID
	title       Title
	description Description
	startTime   time.Time
	endTime     time.Time
	color       Color
	createdAt   time.Time
	updatedAt   time.Time
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

func (e event) Color() Color {
	return e.color
}

func (e event) CreatedAt() time.Time {
	return e.createdAt
}

func (e event) UpdatedAt() time.Time {
	return e.updatedAt
}

func (e *event) Update(title Title, description Description, startTime, endTime time.Time, color Color) {
	e.title = title
	e.description = description
	e.startTime = startTime
	e.endTime = endTime
	e.color = color
	e.updatedAt = time.Now()
}

func NewEvent(
	id EventID,
	userID user.UserID,
	title Title,
	description Description,
	startTime time.Time,
	endTime time.Time,
	color Color,
	createdAt time.Time,
	updatedAt time.Time,
) Event {
	return &event{
		id:          id,
		userID:      userID,
		title:       title,
		description: description,
		startTime:   startTime,
		endTime:     endTime,
		color:       color,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}
