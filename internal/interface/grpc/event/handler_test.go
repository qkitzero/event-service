package event

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/timestamppb"

	eventv1 "github.com/qkitzero/event-service/gen/go/event/v1"
	"github.com/qkitzero/event-service/internal/domain/event"
	mocksappevent "github.com/qkitzero/event-service/mocks/application/event"
	mocksevent "github.com/qkitzero/event-service/mocks/domain/event"
)

func TestCreateEvent(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		success        bool
		ctx            context.Context
		title          string
		description    string
		startTime      *timestamppb.Timestamp
		endTime        *timestamppb.Timestamp
		color          *string
		createEventErr error
	}{
		{"success create event", true, context.Background(), "title", "description", timestamppb.Now(), timestamppb.Now(), func(s string) *string { return &s }("#FFFFFF"), nil},
		{"failure create event error", false, context.Background(), "title", "description", timestamppb.Now(), timestamppb.Now(), func(s string) *string { return &s }("#FFFFFF"), fmt.Errorf("create event error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockEventUsecase := mocksappevent.NewMockEventUsecase(ctrl)
			mockEvent := mocksevent.NewMockEvent(ctrl)
			mockEventUsecase.EXPECT().CreateEvent(tt.ctx, tt.title, tt.description, tt.startTime, tt.endTime, *tt.color).Return(mockEvent, tt.createEventErr).AnyTimes()
			mockEvent.EXPECT().ID().Return(event.NewEventID()).AnyTimes()
			mockEvent.EXPECT().Title().Return(event.Title(tt.title)).AnyTimes()
			mockEvent.EXPECT().Description().Return(event.Description(tt.description)).AnyTimes()
			mockEvent.EXPECT().StartTime().Return(tt.startTime.AsTime()).AnyTimes()
			mockEvent.EXPECT().EndTime().Return(tt.endTime.AsTime()).AnyTimes()
			mockEvent.EXPECT().Color().Return(event.Color(*tt.color)).AnyTimes()

			eventHandler := NewEventHandler(mockEventUsecase)

			req := &eventv1.CreateEventRequest{
				Title:       tt.title,
				Description: tt.description,
				StartTime:   tt.startTime,
				EndTime:     tt.endTime,
				Color:       tt.color,
			}

			_, err := eventHandler.CreateEvent(tt.ctx, req)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}
func TestUpdateEvent(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		success        bool
		ctx            context.Context
		title          string
		id             string
		description    string
		startTime      *timestamppb.Timestamp
		endTime        *timestamppb.Timestamp
		color          string
		updateEventErr error
	}{
		{"success update event", true, context.Background(), "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", "title", "description", timestamppb.Now(), timestamppb.Now(), "#FFFFFF", nil},
		{"failure update event error", false, context.Background(), "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", "title", "description", timestamppb.Now(), timestamppb.Now(), "#FFFFFF", fmt.Errorf("update event error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockEventUsecase := mocksappevent.NewMockEventUsecase(ctrl)
			mockEvent := mocksevent.NewMockEvent(ctrl)
			mockEventUsecase.EXPECT().UpdateEvent(tt.ctx, tt.id, tt.title, tt.description, tt.startTime, tt.endTime, tt.color).Return(mockEvent, tt.updateEventErr).AnyTimes()
			mockEvent.EXPECT().ID().Return(event.NewEventID()).AnyTimes()
			mockEvent.EXPECT().Title().Return(event.Title(tt.title)).AnyTimes()
			mockEvent.EXPECT().Description().Return(event.Description(tt.description)).AnyTimes()
			mockEvent.EXPECT().StartTime().Return(tt.startTime.AsTime()).AnyTimes()
			mockEvent.EXPECT().EndTime().Return(tt.endTime.AsTime()).AnyTimes()
			mockEvent.EXPECT().Color().Return(event.Color(tt.color)).AnyTimes()

			eventHandler := NewEventHandler(mockEventUsecase)

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

			_, err := eventHandler.UpdateEvent(tt.ctx, req)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}

func TestGetEvent(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		success     bool
		ctx         context.Context
		id          string
		getEventErr error
	}{
		{"success get event", true, context.Background(), "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", nil},
		{"failure get event error", false, context.Background(), "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", fmt.Errorf("get event error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockEventUsecase := mocksappevent.NewMockEventUsecase(ctrl)
			mockEvent := mocksevent.NewMockEvent(ctrl)
			mockEventUsecase.EXPECT().GetEvent(tt.ctx, tt.id).Return(mockEvent, tt.getEventErr).AnyTimes()
			mockEvent.EXPECT().ID().Return(event.NewEventID()).AnyTimes()
			mockEvent.EXPECT().Title().Return(event.Title("title")).AnyTimes()
			mockEvent.EXPECT().Description().Return(event.Description("description")).AnyTimes()
			mockEvent.EXPECT().StartTime().Return(time.Now()).AnyTimes()
			mockEvent.EXPECT().EndTime().Return(time.Now()).AnyTimes()
			mockEvent.EXPECT().Color().Return(event.Color("#FFFFFF")).AnyTimes()

			eventHandler := NewEventHandler(mockEventUsecase)

			req := &eventv1.GetEventRequest{
				Id: tt.id,
			}

			_, err := eventHandler.GetEvent(tt.ctx, req)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}

func TestListEvents(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		success       bool
		ctx           context.Context
		listEventsErr error
	}{
		{"success list events", true, context.Background(), nil},
		{"failure list events error", false, context.Background(), fmt.Errorf("list events error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockEventUsecase := mocksappevent.NewMockEventUsecase(ctrl)
			mockEvent := mocksevent.NewMockEvent(ctrl)
			mockEventUsecase.EXPECT().ListEvents(tt.ctx).Return([]event.Event{mockEvent}, tt.listEventsErr).AnyTimes()
			mockEvent.EXPECT().ID().Return(event.NewEventID()).AnyTimes()
			mockEvent.EXPECT().Title().Return(event.Title("title")).AnyTimes()
			mockEvent.EXPECT().Description().Return(event.Description("description")).AnyTimes()
			mockEvent.EXPECT().StartTime().Return(time.Now()).AnyTimes()
			mockEvent.EXPECT().EndTime().Return(time.Now()).AnyTimes()
			mockEvent.EXPECT().Color().Return(event.Color("#FFFFFF")).AnyTimes()

			eventHandler := NewEventHandler(mockEventUsecase)

			req := &eventv1.ListEventsRequest{}

			_, err := eventHandler.ListEvents(tt.ctx, req)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}

func TestDeleteEvent(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		success        bool
		ctx            context.Context
		id             string
		deleteEventErr error
	}{
		{"success delete event", true, context.Background(), "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", nil},
		{"failure delete event error", false, context.Background(), "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", fmt.Errorf("delete event error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockEventUsecase := mocksappevent.NewMockEventUsecase(ctrl)
			mockEventUsecase.EXPECT().DeleteEvent(tt.ctx, tt.id).Return(tt.deleteEventErr).AnyTimes()

			eventHandler := NewEventHandler(mockEventUsecase)

			req := &eventv1.DeleteEventRequest{
				Id: tt.id,
			}

			_, err := eventHandler.DeleteEvent(tt.ctx, req)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}
