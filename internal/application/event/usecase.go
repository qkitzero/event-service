package event

import (
	"context"
	"time"

	"github.com/qkitzero/event-service/internal/application/user"
	"github.com/qkitzero/event-service/internal/domain/event"
	domainuser "github.com/qkitzero/event-service/internal/domain/user"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type EventUsecase interface {
	CreateEvent(ctx context.Context, title, description string, startTime, endTime *timestamppb.Timestamp, color string) (event.Event, error)
	UpdateEvent(ctx context.Context, eventID, title, description string, startTime, endTime *timestamppb.Timestamp, color string) (event.Event, error)
	GetEvent(ctx context.Context, eventID string) (event.Event, error)
	ListEvents(ctx context.Context) ([]event.Event, error)
	DeleteEvent(ctx context.Context, eventID string) error
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

func (s *eventUsecase) CreateEvent(ctx context.Context, title, description string, startTime, endTime *timestamppb.Timestamp, color string) (event.Event, error) {
	userID, err := s.userService.GetUser(ctx)
	if err != nil {
		return nil, err
	}

	newUserID, err := domainuser.NewUserIDFromString(userID)
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
		return nil, event.ErrStartTimeRequired
	}
	newStartTime := startTime.AsTime()

	if endTime == nil {
		return nil, event.ErrEndTimeRequired
	}
	newEndTime := endTime.AsTime()

	newColor, err := event.NewColor(color)
	if err != nil {
		return nil, err
	}

	newEvent := event.NewEvent(event.NewEventID(), newUserID, newTitle, newDescription, newStartTime, newEndTime, newColor, time.Now(), time.Now())

	if err := s.eventRepo.Create(newEvent); err != nil {
		return nil, err
	}

	return newEvent, nil
}

func (s *eventUsecase) UpdateEvent(ctx context.Context, eventID, title, description string, startTime, endTime *timestamppb.Timestamp, color string) (event.Event, error) {
	userID, err := s.userService.GetUser(ctx)
	if err != nil {
		return nil, err
	}

	id, err := event.NewEventIDFromString(eventID)
	if err != nil {
		return nil, err
	}

	foundEvent, err := s.eventRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if foundEvent.UserID().String() != userID {
		return nil, event.ErrPermissionDenied
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

	if err := s.eventRepo.Update(foundEvent); err != nil {
		return nil, err
	}

	return foundEvent, nil
}

func (s *eventUsecase) GetEvent(ctx context.Context, eventID string) (event.Event, error) {
	userID, err := s.userService.GetUser(ctx)
	if err != nil {
		return nil, err
	}

	id, err := event.NewEventIDFromString(eventID)
	if err != nil {
		return nil, err
	}

	foundEvent, err := s.eventRepo.FindByID(id)
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

func (s *eventUsecase) DeleteEvent(ctx context.Context, eventID string) error {
	userID, err := s.userService.GetUser(ctx)
	if err != nil {
		return err
	}

	id, err := event.NewEventIDFromString(eventID)
	if err != nil {
		return err
	}

	foundEvent, err := s.eventRepo.FindByID(id)
	if err != nil {
		return err
	}

	if foundEvent.UserID().String() != userID {
		return event.ErrPermissionDenied
	}

	if err := s.eventRepo.Delete(id); err != nil {
		return err
	}

	return nil
}
