package event

import (
	"time"

	"github.com/qkitzero/event-service/internal/domain/event"
	"github.com/qkitzero/event-service/internal/domain/user"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type EventUsecase interface {
	CreateEvent(userIDStr, titleStr, descriptionStr string, startTime, endTime time.Time) (event.Event, error)
	UpdateEvent(eventID, title, description string, startTime, endTime *timestamppb.Timestamp) (event.Event, error)
	GetEvent(eventIDStr string) (event.Event, error)
	ListEvents(userIDStr string) ([]event.Event, error)
	DeleteEvent(eventIDStr string) error
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

	e := event.NewEvent(event.NewEventID(), userID, title, description, startTime, endTime, time.Now(), time.Now())

	if err := s.repo.Create(e); err != nil {
		return nil, err
	}

	return e, nil
}

func (s *eventUsecase) UpdateEvent(eventID, title, description string, startTime, endTime *timestamppb.Timestamp) (event.Event, error) {
	id, err := event.NewEventIDFromString(eventID)
	if err != nil {
		return nil, err
	}

	foundEvent, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	newTitle, err := event.NewTitle(title)
	if err != nil {
		return nil, err
	}

	newDescription, err := event.NewDescription(description)
	if err != nil {
		return nil, err
	}

	newStartTime := foundEvent.StartTime()
	if startTime != nil {
		newStartTime = startTime.AsTime()
	}

	newEndTime := foundEvent.EndTime()
	if endTime != nil {
		newEndTime = endTime.AsTime()
	}

	foundEvent.Update(newTitle, newDescription, newStartTime, newEndTime)

	if err := s.repo.Update(foundEvent); err != nil {
		return nil, err
	}

	return foundEvent, nil
}

func (s *eventUsecase) GetEvent(eventIDStr string) (event.Event, error) {
	eventID, err := event.NewEventIDFromString(eventIDStr)
	if err != nil {
		return nil, err
	}

	e, err := s.repo.FindByID(eventID)
	if err != nil {
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

func (s *eventUsecase) DeleteEvent(eventIDStr string) error {
	eventID, err := event.NewEventIDFromString(eventIDStr)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(eventID); err != nil {
		return err
	}

	return nil
}
