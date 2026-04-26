package event

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	eventv1 "github.com/qkitzero/event-service/gen/go/event/v1"
	"github.com/qkitzero/event-service/internal/domain/event"
	mocksappevent "github.com/qkitzero/event-service/mocks/application/event"
	mocksevent "github.com/qkitzero/event-service/mocks/domain/event"
)

func mockEventSample(ctrl *gomock.Controller) *mocksevent.MockEvent {
	m := mocksevent.NewMockEvent(ctrl)
	m.EXPECT().ID().Return(event.NewEventID()).AnyTimes()
	m.EXPECT().Title().Return(event.Title("title")).AnyTimes()
	m.EXPECT().Description().Return(event.Description("description")).AnyTimes()
	m.EXPECT().StartTime().Return(time.Now()).AnyTimes()
	m.EXPECT().EndTime().Return(time.Now()).AnyTimes()
	m.EXPECT().Color().Return(event.Color("#FFFFFF")).AnyTimes()
	return m
}

func ptr(s string) *string { return &s }

func TestCreateEvent(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		ctx            context.Context
		title          string
		description    string
		startTime      *timestamppb.Timestamp
		endTime        *timestamppb.Timestamp
		color          *string
		callUsecase    bool
		createEventErr error
		wantCode       codes.Code
	}{
		{"success create event", context.Background(), "title", "description", timestamppb.Now(), timestamppb.Now(), ptr("#FFFFFF"), true, nil, codes.OK},
		{"failure empty title", context.Background(), "", "description", timestamppb.Now(), timestamppb.Now(), ptr("#FFFFFF"), false, nil, codes.InvalidArgument},
		{"failure empty description", context.Background(), "title", "", timestamppb.Now(), timestamppb.Now(), ptr("#FFFFFF"), false, nil, codes.InvalidArgument},
		{"failure nil start time", context.Background(), "title", "description", nil, timestamppb.Now(), ptr("#FFFFFF"), false, nil, codes.InvalidArgument},
		{"failure nil end time", context.Background(), "title", "description", timestamppb.Now(), nil, ptr("#FFFFFF"), false, nil, codes.InvalidArgument},
		{"failure invalid color", context.Background(), "title", "description", timestamppb.Now(), timestamppb.Now(), ptr("red"), false, nil, codes.InvalidArgument},
		{"failure event not found", context.Background(), "title", "description", timestamppb.Now(), timestamppb.Now(), ptr("#FFFFFF"), true, event.ErrEventNotFound, codes.NotFound},
		{"failure permission denied", context.Background(), "title", "description", timestamppb.Now(), timestamppb.Now(), ptr("#FFFFFF"), true, event.ErrPermissionDenied, codes.PermissionDenied},
		{"failure usecase error", context.Background(), "title", "description", timestamppb.Now(), timestamppb.Now(), ptr("#FFFFFF"), true, fmt.Errorf("create event error"), codes.Internal},
		{"failure status preserved", context.Background(), "title", "description", timestamppb.Now(), timestamppb.Now(), ptr("#FFFFFF"), true, status.Error(codes.Unauthenticated, "auth"), codes.Unauthenticated},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocksappevent.NewMockEventUsecase(ctrl)
			if tt.callUsecase {
				mockUsecase.EXPECT().CreateEvent(tt.ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockEventSample(ctrl), tt.createEventErr).Times(1)
			}

			handler := NewEventHandler(mockUsecase)

			req := &eventv1.CreateEventRequest{
				Title:       tt.title,
				Description: tt.description,
				StartTime:   tt.startTime,
				EndTime:     tt.endTime,
				Color:       tt.color,
			}

			_, err := handler.CreateEvent(tt.ctx, req)
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
		})
	}
}

func TestUpdateEvent(t *testing.T) {
	t.Parallel()
	validID := "fe8c2263-bbac-4bb9-a41d-b04f5afc4425"
	tests := []struct {
		name           string
		ctx            context.Context
		id             string
		title          string
		description    string
		startTime      *timestamppb.Timestamp
		endTime        *timestamppb.Timestamp
		color          string
		callUsecase    bool
		updateEventErr error
		wantCode       codes.Code
	}{
		{"success update event", context.Background(), validID, "title", "description", timestamppb.Now(), timestamppb.Now(), "#FFFFFF", true, nil, codes.OK},
		{"failure invalid event id", context.Background(), "not-uuid", "title", "description", timestamppb.Now(), timestamppb.Now(), "#FFFFFF", false, nil, codes.InvalidArgument},
		{"failure empty title", context.Background(), validID, "", "description", timestamppb.Now(), timestamppb.Now(), "#FFFFFF", false, nil, codes.InvalidArgument},
		{"failure empty description", context.Background(), validID, "title", "", timestamppb.Now(), timestamppb.Now(), "#FFFFFF", false, nil, codes.InvalidArgument},
		{"failure invalid color", context.Background(), validID, "title", "description", timestamppb.Now(), timestamppb.Now(), "red", false, nil, codes.InvalidArgument},
		{"failure event not found", context.Background(), validID, "title", "description", timestamppb.Now(), timestamppb.Now(), "#FFFFFF", true, event.ErrEventNotFound, codes.NotFound},
		{"failure permission denied", context.Background(), validID, "title", "description", timestamppb.Now(), timestamppb.Now(), "#FFFFFF", true, event.ErrPermissionDenied, codes.PermissionDenied},
		{"failure usecase error", context.Background(), validID, "title", "description", timestamppb.Now(), timestamppb.Now(), "#FFFFFF", true, fmt.Errorf("update event error"), codes.Internal},
		{"failure status preserved", context.Background(), validID, "title", "description", timestamppb.Now(), timestamppb.Now(), "#FFFFFF", true, status.Error(codes.Unauthenticated, "auth"), codes.Unauthenticated},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocksappevent.NewMockEventUsecase(ctrl)
			if tt.callUsecase {
				mockUsecase.EXPECT().UpdateEvent(tt.ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockEventSample(ctrl), tt.updateEventErr).Times(1)
			}

			handler := NewEventHandler(mockUsecase)

			req := &eventv1.UpdateEventRequest{
				Event: &eventv1.Event{
					Id:          tt.id,
					Title:       tt.title,
					Description: tt.description,
					StartTime:   tt.startTime,
					EndTime:     tt.endTime,
					Color:       tt.color,
				},
			}

			_, err := handler.UpdateEvent(tt.ctx, req)
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
		})
	}
}

func TestGetEvent(t *testing.T) {
	t.Parallel()
	validID := "fe8c2263-bbac-4bb9-a41d-b04f5afc4425"
	tests := []struct {
		name        string
		ctx         context.Context
		id          string
		callUsecase bool
		getEventErr error
		wantCode    codes.Code
	}{
		{"success get event", context.Background(), validID, true, nil, codes.OK},
		{"failure invalid id", context.Background(), "not-uuid", false, nil, codes.InvalidArgument},
		{"failure event not found", context.Background(), validID, true, event.ErrEventNotFound, codes.NotFound},
		{"failure permission denied", context.Background(), validID, true, event.ErrPermissionDenied, codes.PermissionDenied},
		{"failure usecase error", context.Background(), validID, true, fmt.Errorf("get event error"), codes.Internal},
		{"failure status preserved", context.Background(), validID, true, status.Error(codes.Unauthenticated, "auth"), codes.Unauthenticated},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocksappevent.NewMockEventUsecase(ctrl)
			if tt.callUsecase {
				mockUsecase.EXPECT().GetEvent(tt.ctx, gomock.Any()).Return(mockEventSample(ctrl), tt.getEventErr).Times(1)
			}

			handler := NewEventHandler(mockUsecase)

			_, err := handler.GetEvent(tt.ctx, &eventv1.GetEventRequest{Id: tt.id})
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
		})
	}
}

func TestListEvents(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		ctx           context.Context
		listEventsErr error
		wantCode      codes.Code
	}{
		{"success list events", context.Background(), nil, codes.OK},
		{"failure list events error", context.Background(), fmt.Errorf("list events error"), codes.Internal},
		{"failure status preserved", context.Background(), status.Error(codes.Unauthenticated, "auth"), codes.Unauthenticated},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocksappevent.NewMockEventUsecase(ctrl)
			mockUsecase.EXPECT().ListEvents(tt.ctx).Return([]event.Event{mockEventSample(ctrl)}, tt.listEventsErr).Times(1)

			handler := NewEventHandler(mockUsecase)

			_, err := handler.ListEvents(tt.ctx, &eventv1.ListEventsRequest{})
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
		})
	}
}

func TestDeleteEvent(t *testing.T) {
	t.Parallel()
	validID := "fe8c2263-bbac-4bb9-a41d-b04f5afc4425"
	tests := []struct {
		name           string
		ctx            context.Context
		id             string
		callUsecase    bool
		deleteEventErr error
		wantCode       codes.Code
	}{
		{"success delete event", context.Background(), validID, true, nil, codes.OK},
		{"failure invalid id", context.Background(), "not-uuid", false, nil, codes.InvalidArgument},
		{"failure event not found", context.Background(), validID, true, event.ErrEventNotFound, codes.NotFound},
		{"failure permission denied", context.Background(), validID, true, event.ErrPermissionDenied, codes.PermissionDenied},
		{"failure usecase error", context.Background(), validID, true, fmt.Errorf("delete event error"), codes.Internal},
		{"failure status preserved", context.Background(), validID, true, status.Error(codes.Unauthenticated, "auth"), codes.Unauthenticated},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocksappevent.NewMockEventUsecase(ctrl)
			if tt.callUsecase {
				mockUsecase.EXPECT().DeleteEvent(tt.ctx, gomock.Any()).Return(tt.deleteEventErr).Times(1)
			}

			handler := NewEventHandler(mockUsecase)

			_, err := handler.DeleteEvent(tt.ctx, &eventv1.DeleteEventRequest{Id: tt.id})
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
		})
	}
}
