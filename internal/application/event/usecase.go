package event

import (
	"time"

	"github.com/qkitzero/event-service/internal/domain/event"
	"github.com/qkitzero/event-service/internal/domain/user"
)

type EventUsecase interface {
	CreateEvent(userIDStr, titleStr, descriptionStr string, startTime, endTime time.Time) (event.Event, error)
	ListEvents(userIDStr string) ([]event.Event, error)
}

type eventUsecase struct {
	repo event.EventRepository
}

func NewEventUsecase(repo event.EventRepository) EventUsecase {
	return &eventUsecase{repo: repo}
}

func (s *eventUsecase) CreateEvent(userIDStr, titleStr, descriptionStr string, startTime, endTime time.Time) (event.Event, error) {
	userID, err := user.NewUserIDFromString(userIDStr)
	if err != nil {
		return nil, err
	}

	title, err := event.NewTitle(titleStr)
	if err != nil {
		return nil, err
	}

	description, err := event.NewDescription(descriptionStr)
	if err != nil {
		return nil, err
	}

	e := event.NewEvent(event.NewEventID(), userID, title, description, startTime, endTime, time.Now())

	if err := s.repo.Create(e); err != nil {
		return nil, err
	}

	return e, nil
}

func (s *eventUsecase) ListEvents(userIDStr string) ([]event.Event, error) {
	userID, err := user.NewUserIDFromString(userIDStr)
	if err != nil {
		return nil, err
	}

	events, err := s.repo.ListByUserID(userID)
	if err != nil {
		return nil, err
	}

	return events, nil
}
