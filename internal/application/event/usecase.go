package event

import (
	"fmt"
	"time"

	"github.com/qkitzero/event-service/internal/domain/event"
	"github.com/qkitzero/event-service/internal/domain/user"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type EventUsecase interface {
	CreateEvent(userID, title, description string, startTime, endTime *timestamppb.Timestamp, color string) (event.Event, error)
	UpdateEvent(eventID, title, description string, startTime, endTime *timestamppb.Timestamp, color string) (event.Event, error)
	GetEvent(eventID string) (event.Event, error)
	ListEvents(userID string) ([]event.Event, error)
	DeleteEvent(eventID string) error
}

type eventUsecase struct {
	repo event.EventRepository
}

func NewEventUsecase(repo event.EventRepository) EventUsecase {
	return &eventUsecase{repo: repo}
}

func (s *eventUsecase) CreateEvent(userID, title, description string, startTime, endTime *timestamppb.Timestamp, color string) (event.Event, error) {
	newUserID, err := user.NewUserIDFromString(userID)
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

	if startTime == nil {
		return nil, fmt.Errorf("start time is required")
	}
	newStartTime := startTime.AsTime()

	if endTime == nil {
		return nil, fmt.Errorf("end time is required")
	}
	newEndTime := endTime.AsTime()

	newColor, err := event.NewColor(color)
	if err != nil {
		return nil, err
	}

	newEvent := event.NewEvent(event.NewEventID(), newUserID, newTitle, newDescription, newStartTime, newEndTime, newColor, time.Now(), time.Now())

	if err := s.repo.Create(newEvent); err != nil {
		return nil, err
	}

	return newEvent, nil
}

func (s *eventUsecase) UpdateEvent(eventID, title, description string, startTime, endTime *timestamppb.Timestamp, color string) (event.Event, error) {
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

	newColor, err := event.NewColor(color)
	if err != nil {
		return nil, err
	}

	foundEvent.Update(newTitle, newDescription, newStartTime, newEndTime, newColor)

	if err := s.repo.Update(foundEvent); err != nil {
		return nil, err
	}

	return foundEvent, nil
}

func (s *eventUsecase) GetEvent(eventID string) (event.Event, error) {
	id, err := event.NewEventIDFromString(eventID)
	if err != nil {
		return nil, err
	}

	foundEvent, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return foundEvent, nil
}

func (s *eventUsecase) ListEvents(userID string) ([]event.Event, error) {
	uid, err := user.NewUserIDFromString(userID)
	if err != nil {
		return nil, err
	}

	events, err := s.repo.FindAllByUserID(uid)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (s *eventUsecase) DeleteEvent(eventID string) error {
	id, err := event.NewEventIDFromString(eventID)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(id); err != nil {
		return err
	}

	return nil
}
