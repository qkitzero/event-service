package event

import (
	"context"
	"errors"
	"log"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	eventv1 "github.com/qkitzero/event-service/gen/go/event/v1"
	appevent "github.com/qkitzero/event-service/internal/application/event"
	domainevent "github.com/qkitzero/event-service/internal/domain/event"
)

type EventHandler struct {
	eventv1.UnimplementedEventServiceServer
	eventUsecase appevent.EventUsecase
}

func NewEventHandler(eventUsecase appevent.EventUsecase) *EventHandler {
	return &EventHandler{
		eventUsecase: eventUsecase,
	}
}

func (h *EventHandler) CreateEvent(ctx context.Context, req *eventv1.CreateEventRequest) (*eventv1.CreateEventResponse, error) {
	title, err := domainevent.NewTitle(req.GetTitle())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	description, err := domainevent.NewDescription(req.GetDescription())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if req.GetStartTime() == nil {
		return nil, status.Error(codes.InvalidArgument, domainevent.ErrStartTimeRequired.Error())
	}
	if req.GetEndTime() == nil {
		return nil, status.Error(codes.InvalidArgument, domainevent.ErrEndTimeRequired.Error())
	}
	color, err := domainevent.NewColor(req.GetColor())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	created, err := h.eventUsecase.CreateEvent(ctx, title, description, req.GetStartTime().AsTime(), req.GetEndTime().AsTime(), color)
	if err != nil {
		if _, ok := status.FromError(err); ok {
			return nil, err
		}
		if errors.Is(err, domainevent.ErrEventNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		if errors.Is(err, domainevent.ErrPermissionDenied) {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}
		log.Printf("CreateEvent: internal error: %v", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &eventv1.CreateEventResponse{
		Event: &eventv1.Event{
			Id:          created.ID().String(),
			Title:       created.Title().String(),
			Description: created.Description().String(),
			StartTime:   timestamppb.New(created.StartTime()),
			EndTime:     timestamppb.New(created.EndTime()),
			Color:       created.Color().String(),
		},
	}, nil
}

func (h *EventHandler) UpdateEvent(ctx context.Context, req *eventv1.UpdateEventRequest) (*eventv1.UpdateEventResponse, error) {
	src := req.GetEvent()

	eventID, err := domainevent.NewEventIDFromString(src.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	title, err := domainevent.NewTitle(src.GetTitle())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	description, err := domainevent.NewDescription(src.GetDescription())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	color, err := domainevent.NewColor(src.GetColor())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var startTime, endTime *time.Time
	if src.GetStartTime() != nil {
		t := src.GetStartTime().AsTime()
		startTime = &t
	}
	if src.GetEndTime() != nil {
		t := src.GetEndTime().AsTime()
		endTime = &t
	}

	updated, err := h.eventUsecase.UpdateEvent(ctx, eventID, title, description, startTime, endTime, color)
	if err != nil {
		if _, ok := status.FromError(err); ok {
			return nil, err
		}
		if errors.Is(err, domainevent.ErrEventNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		if errors.Is(err, domainevent.ErrPermissionDenied) {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}
		log.Printf("UpdateEvent: internal error: %v", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &eventv1.UpdateEventResponse{
		Event: &eventv1.Event{
			Id:          updated.ID().String(),
			Title:       updated.Title().String(),
			Description: updated.Description().String(),
			StartTime:   timestamppb.New(updated.StartTime()),
			EndTime:     timestamppb.New(updated.EndTime()),
			Color:       updated.Color().String(),
		},
	}, nil
}

func (h *EventHandler) GetEvent(ctx context.Context, req *eventv1.GetEventRequest) (*eventv1.GetEventResponse, error) {
	eventID, err := domainevent.NewEventIDFromString(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	got, err := h.eventUsecase.GetEvent(ctx, eventID)
	if err != nil {
		if _, ok := status.FromError(err); ok {
			return nil, err
		}
		if errors.Is(err, domainevent.ErrEventNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		if errors.Is(err, domainevent.ErrPermissionDenied) {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}
		log.Printf("GetEvent: internal error: %v", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &eventv1.GetEventResponse{
		Event: &eventv1.Event{
			Id:          got.ID().String(),
			Title:       got.Title().String(),
			Description: got.Description().String(),
			StartTime:   timestamppb.New(got.StartTime()),
			EndTime:     timestamppb.New(got.EndTime()),
			Color:       got.Color().String(),
		},
	}, nil
}

func (h *EventHandler) ListEvents(ctx context.Context, req *eventv1.ListEventsRequest) (*eventv1.ListEventsResponse, error) {
	events, err := h.eventUsecase.ListEvents(ctx)
	if err != nil {
		if _, ok := status.FromError(err); ok {
			return nil, err
		}
		log.Printf("ListEvents: internal error: %v", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	pbEvents := make([]*eventv1.Event, 0, len(events))
	for _, e := range events {
		pbEvents = append(pbEvents, &eventv1.Event{
			Id:          e.ID().String(),
			Title:       e.Title().String(),
			Description: e.Description().String(),
			StartTime:   timestamppb.New(e.StartTime()),
			EndTime:     timestamppb.New(e.EndTime()),
			Color:       e.Color().String(),
		})
	}

	return &eventv1.ListEventsResponse{
		Events: pbEvents,
	}, nil
}

func (h *EventHandler) DeleteEvent(ctx context.Context, req *eventv1.DeleteEventRequest) (*eventv1.DeleteEventResponse, error) {
	eventID, err := domainevent.NewEventIDFromString(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := h.eventUsecase.DeleteEvent(ctx, eventID); err != nil {
		if _, ok := status.FromError(err); ok {
			return nil, err
		}
		if errors.Is(err, domainevent.ErrEventNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		if errors.Is(err, domainevent.ErrPermissionDenied) {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}
		log.Printf("DeleteEvent: internal error: %v", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &eventv1.DeleteEventResponse{}, nil
}
