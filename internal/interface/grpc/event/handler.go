package event

import (
	"context"

	eventv1 "github.com/qkitzero/event-service/gen/go/event/v1"
	appauth "github.com/qkitzero/event-service/internal/application/auth"
	appevent "github.com/qkitzero/event-service/internal/application/event"
	appuser "github.com/qkitzero/event-service/internal/application/user"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type EventHandler struct {
	eventv1.UnimplementedEventServiceServer
	authUsecase  appauth.AuthUsecase
	userUsecase  appuser.UserUsecase
	eventUsecase appevent.EventUsecase
}

func NewEventHandler(
	authUsecase appauth.AuthUsecase,
	userUsecase appuser.UserUsecase,
	eventUsecase appevent.EventUsecase,
) *EventHandler {
	return &EventHandler{
		authUsecase:  authUsecase,
		userUsecase:  userUsecase,
		eventUsecase: eventUsecase,
	}
}

func (h *EventHandler) CreateEvent(ctx context.Context, req *eventv1.CreateEventRequest) (*eventv1.CreateEventResponse, error) {
	userID, err := h.userUsecase.GetUser(ctx)
	if err != nil {
		return nil, err
	}

	event, err := h.eventUsecase.CreateEvent(userID, req.GetTitle(), req.GetDescription(), req.GetStartTime(), req.GetEndTime())
	if err != nil {
		return nil, err
	}

	return &eventv1.CreateEventResponse{
		Event: &eventv1.Event{
			Id:          event.ID().String(),
			Title:       event.Title().String(),
			Description: event.Description().String(),
			StartTime:   timestamppb.New(event.StartTime()),
			EndTime:     timestamppb.New(event.EndTime()),
		},
	}, nil
}

func (h *EventHandler) UpdateEvent(ctx context.Context, req *eventv1.UpdateEventRequest) (*eventv1.UpdateEventResponse, error) {
	_, err := h.authUsecase.VerifyToken(ctx)
	if err != nil {
		return nil, err
	}

	event, err := h.eventUsecase.UpdateEvent(req.GetEvent().GetId(), req.GetEvent().GetTitle(), req.GetEvent().GetDescription(), req.GetEvent().GetStartTime(), req.GetEvent().GetEndTime())
	if err != nil {
		return nil, err
	}

	return &eventv1.UpdateEventResponse{
		Event: &eventv1.Event{
			Id:          event.ID().String(),
			Title:       event.Title().String(),
			Description: event.Description().String(),
			StartTime:   timestamppb.New(event.StartTime()),
			EndTime:     timestamppb.New(event.EndTime()),
		},
	}, nil
}

func (h *EventHandler) GetEvent(ctx context.Context, req *eventv1.GetEventRequest) (*eventv1.GetEventResponse, error) {
	_, err := h.authUsecase.VerifyToken(ctx)
	if err != nil {
		return nil, err
	}

	event, err := h.eventUsecase.GetEvent(req.GetId())
	if err != nil {
		return nil, err
	}

	return &eventv1.GetEventResponse{
		Event: &eventv1.Event{
			Id:          event.ID().String(),
			Title:       event.Title().String(),
			Description: event.Description().String(),
			StartTime:   timestamppb.New(event.StartTime()),
			EndTime:     timestamppb.New(event.EndTime()),
		},
	}, nil
}

func (h *EventHandler) ListEvents(ctx context.Context, req *eventv1.ListEventsRequest) (*eventv1.ListEventsResponse, error) {
	userID, err := h.userUsecase.GetUser(ctx)
	if err != nil {
		return nil, err
	}

	events, err := h.eventUsecase.ListEvents(userID)
	if err != nil {
		return nil, err
	}

	var pbEvents []*eventv1.Event
	for _, event := range events {
		pbEvent := &eventv1.Event{
			Id:          event.ID().String(),
			Title:       event.Title().String(),
			Description: event.Description().String(),
			StartTime:   timestamppb.New(event.StartTime()),
			EndTime:     timestamppb.New(event.EndTime()),
		}
		pbEvents = append(pbEvents, pbEvent)
	}

	return &eventv1.ListEventsResponse{
		Events: pbEvents,
	}, nil
}

func (h *EventHandler) DeleteEvent(ctx context.Context, req *eventv1.DeleteEventRequest) (*eventv1.DeleteEventResponse, error) {
	_, err := h.authUsecase.VerifyToken(ctx)
	if err != nil {
		return nil, err
	}

	if err := h.eventUsecase.DeleteEvent(req.GetId()); err != nil {
		return nil, err
	}

	return &eventv1.DeleteEventResponse{}, nil
}
