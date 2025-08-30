package event

import (
	"context"

	eventv1 "github.com/qkitzero/event-service/gen/go/event/v1"
	appevent "github.com/qkitzero/event-service/internal/application/event"
	appuser "github.com/qkitzero/event-service/internal/application/user"
)

type EventHandler struct {
	eventv1.UnimplementedEventServiceServer
	userService  appuser.UserUsecase
	eventUsecase appevent.EventUsecase
}

func NewEventHandler(
	userService appuser.UserUsecase,
	eventUsecase appevent.EventUsecase,
) *EventHandler {
	return &EventHandler{
		userService:  userService,
		eventUsecase: eventUsecase,
	}
}

func (h *EventHandler) CreateEvent(ctx context.Context, req *eventv1.CreateEventRequest) (*eventv1.CreateEventResponse, error) {
	userID, err := h.userService.GetUser(ctx)
	if err != nil {
		return nil, err
	}

	event, err := h.eventUsecase.CreateEvent(userID, req.GetTitle(), req.GetDescription(), req.GetStartTime().AsTime(), req.GetEndTime().AsTime())
	if err != nil {
		return nil, err
	}

	return &eventv1.CreateEventResponse{
		EventId: event.ID().String(),
	}, nil
}
