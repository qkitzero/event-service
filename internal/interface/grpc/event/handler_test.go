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
	mocksappuser "github.com/qkitzero/event-service/mocks/application/user"
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
		getUserErr     error
		createEventErr error
	}{
		{"success create event", true, context.Background(), "title", "description", timestamppb.Now(), timestamppb.Now(), nil, nil},
		{"failure get user error", false, context.Background(), "title", "description", timestamppb.Now(), timestamppb.Now(), fmt.Errorf("get user error"), nil},
		{"failure create event error", false, context.Background(), "title", "description", timestamppb.Now(), timestamppb.Now(), nil, fmt.Errorf("create event error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserUsecase := mocksappuser.NewMockUserUsecase(ctrl)
			mockEventUsecase := mocksappevent.NewMockEventUsecase(ctrl)
			mockEvent := mocksevent.NewMockEvent(ctrl)
			mockUserID := "mockUserID"
			mockUserUsecase.EXPECT().GetUser(tt.ctx).Return(mockUserID, tt.getUserErr).AnyTimes()
			mockEventUsecase.EXPECT().CreateEvent(mockUserID, tt.title, tt.description, tt.startTime.AsTime(), tt.endTime.AsTime()).Return(mockEvent, tt.createEventErr).AnyTimes()
			mockEvent.EXPECT().ID().Return(event.NewEventID()).AnyTimes()

			eventHandler := NewEventHandler(mockUserUsecase, mockEventUsecase)

			req := &eventv1.CreateEventRequest{
				Title:       tt.title,
				Description: tt.description,
				StartTime:   tt.startTime,
				EndTime:     tt.endTime,
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

func TestListEvents(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		success       bool
		ctx           context.Context
		getUserErr    error
		listEventsErr error
	}{
		{"success list events", true, context.Background(), nil, nil},
		{"failure get user error", false, context.Background(), fmt.Errorf("get user error"), nil},
		{"failure list events error", false, context.Background(), nil, fmt.Errorf("list events error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserUsecase := mocksappuser.NewMockUserUsecase(ctrl)
			mockEventUsecase := mocksappevent.NewMockEventUsecase(ctrl)
			mockEvent := mocksevent.NewMockEvent(ctrl)
			mockUserID := "mockUserID"
			mockUserUsecase.EXPECT().GetUser(tt.ctx).Return(mockUserID, tt.getUserErr).AnyTimes()
			mockEventUsecase.EXPECT().ListEvents(mockUserID).Return([]event.Event{mockEvent}, tt.listEventsErr).AnyTimes()
			mockEvent.EXPECT().ID().Return(event.NewEventID()).AnyTimes()
			mockEvent.EXPECT().Title().Return(event.Title("title")).AnyTimes()
			mockEvent.EXPECT().Description().Return(event.Description("description")).AnyTimes()
			mockEvent.EXPECT().StartTime().Return(time.Now()).AnyTimes()
			mockEvent.EXPECT().EndTime().Return(time.Now()).AnyTimes()

			eventHandler := NewEventHandler(mockUserUsecase, mockEventUsecase)

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
