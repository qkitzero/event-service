package event

import (
	"context"
	"fmt"
	"testing"

	"go.uber.org/mock/gomock"

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
		getUserErr     error
		createEventErr error
	}{
		{"success create event", true, context.Background(), "title", "description", nil, nil},
		{"failure get user error", false, context.Background(), "title", "description", fmt.Errorf("get user error"), nil},
		{"failure create event error", false, context.Background(), "title", "description", nil, fmt.Errorf("create event error")},
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
			mockEventUsecase.EXPECT().CreateEvent(mockUserID, tt.title, tt.description).Return(mockEvent, tt.createEventErr).AnyTimes()
			mockEventID := event.NewEventID()
			mockEvent.EXPECT().ID().Return(mockEventID).AnyTimes()

			eventHandler := NewEventHandler(mockUserUsecase, mockEventUsecase)

			req := &eventv1.CreateEventRequest{
				Title:       tt.title,
				Description: tt.description,
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
