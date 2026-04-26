package event

import (
	"context"
	"time"

	"github.com/qkitzero/event-service/internal/application/user"
	"github.com/qkitzero/event-service/internal/domain/event"
	domainuser "github.com/qkitzero/event-service/internal/domain/user"
)

type EventUsecase interface {
	CreateEvent(ctx context.Context, title event.Title, description event.Description, startTime, endTime time.Time, color event.Color) (event.Event, error)
	UpdateEvent(ctx context.Context, eventID event.EventID, title event.Title, description event.Description, startTime, endTime *time.Time, color event.Color) (event.Event, error)
	GetEvent(ctx context.Context, eventID event.EventID) (event.Event, error)
	ListEvents(ctx context.Context) ([]event.Event, error)
	DeleteEvent(ctx context.Context, eventID event.EventID) error
}

type eventUsecase struct {
	userService user.UserService
	eventRepo   event.EventRepository
}

func NewEventUsecase(
	userService user.UserService,
	eventRepo event.EventRepository,
) EventUsecase {
	return &eventUsecase{
		userService: userService,
		eventRepo:   eventRepo,
	}
}

func (s *eventUsecase) CreateEvent(ctx context.Context, title event.Title, description event.Description, startTime, endTime time.Time, color event.Color) (event.Event, error) {
	userID, err := s.userService.GetUser(ctx)
	if err != nil {
		return nil, err
	}

	newUserID, err := domainuser.NewUserIDFromString(userID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	newEvent := event.NewEvent(event.NewEventID(), newUserID, title, description, startTime, endTime, color, now, now)

	if err := s.eventRepo.Create(newEvent); err != nil {
		return nil, err
	}

	return newEvent, nil
}

func (s *eventUsecase) UpdateEvent(ctx context.Context, eventID event.EventID, title event.Title, description event.Description, startTime, endTime *time.Time, color event.Color) (event.Event, error) {
	userID, err := s.userService.GetUser(ctx)
	if err != nil {
		return nil, err
	}

	foundEvent, err := s.eventRepo.FindByID(eventID)
	if err != nil {
		return nil, err
	}

	if foundEvent.UserID().String() != userID {
		return nil, event.ErrPermissionDenied
	}

	newStartTime := foundEvent.StartTime()
	if startTime != nil {
		newStartTime = *startTime
	}

	newEndTime := foundEvent.EndTime()
	if endTime != nil {
		newEndTime = *endTime
	}

	foundEvent.Update(title, description, newStartTime, newEndTime, color)

	if err := s.eventRepo.Update(foundEvent); err != nil {
		return nil, err
	}

	return foundEvent, nil
}

func (s *eventUsecase) GetEvent(ctx context.Context, eventID event.EventID) (event.Event, error) {
	userID, err := s.userService.GetUser(ctx)
	if err != nil {
		return nil, err
	}

	foundEvent, err := s.eventRepo.FindByID(eventID)
	if err != nil {
		return nil, err
	}

	if foundEvent.UserID().String() != userID {
		return nil, event.ErrPermissionDenied
	}

	return foundEvent, nil
}

func (s *eventUsecase) ListEvents(ctx context.Context) ([]event.Event, error) {
	userID, err := s.userService.GetUser(ctx)
	if err != nil {
		return nil, err
	}

	uid, err := domainuser.NewUserIDFromString(userID)
	if err != nil {
		return nil, err
	}

	events, err := s.eventRepo.FindAllByUserID(uid)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (s *eventUsecase) DeleteEvent(ctx context.Context, eventID event.EventID) error {
	userID, err := s.userService.GetUser(ctx)
	if err != nil {
		return err
	}

	foundEvent, err := s.eventRepo.FindByID(eventID)
	if err != nil {
		return err
	}

	if foundEvent.UserID().String() != userID {
		return event.ErrPermissionDenied
	}

	if err := s.eventRepo.Delete(eventID); err != nil {
		return err
	}

	return nil
}
