package event

import (
	"context"

	eventv1 "github.com/qkitzero/event-service/gen/go/event/v1"
	appevent "github.com/qkitzero/event-service/internal/application/event"
	appuser "github.com/qkitzero/event-service/internal/application/user"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	event, err := h.eventUsecase.UpdateEvent(req.GetEvent().GetId(), req.GetEvent().GetTitle(), req.GetEvent().GetDescription(), req.GetEvent().GetStartTime().AsTime(), req.GetEvent().GetEndTime().AsTime())
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

func (h *EventHandler) ListEvents(ctx context.Context, req *eventv1.ListEventsRequest) (*eventv1.ListEventsResponse, error) {
	userID, err := h.userService.GetUser(ctx)
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
