package event

import (
	"errors"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	"github.com/qkitzero/event-service/internal/domain/event"
	mocks "github.com/qkitzero/event-service/mocks/domain/event"
)

func TestCreateEvent(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		success     bool
		userID      string
		title       string
		description string
		startTime   time.Time
		endTime     time.Time
		createErr   error
	}{
		{"success create event", true, "6d322c66-bf4d-427a-970c-874f3745f653", "title", "description", time.Now(), time.Now(), nil},
		{"failure empty user id", false, "", "title", "description", time.Now(), time.Now(), nil},
		{"failure empty title", false, "6d322c66-bf4d-427a-970c-874f3745f653", "", "description", time.Now(), time.Now(), nil},
		{"failure empty description", false, "6d322c66-bf4d-427a-970c-874f3745f653", "title", "", time.Now(), time.Now(), nil},
		{"failure create error", false, "6d322c66-bf4d-427a-970c-874f3745f653", "title", "description", time.Now(), time.Now(), errors.New("create error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockEventRepository := mocks.NewMockEventRepository(ctrl)
			mockEventRepository.EXPECT().Create(gomock.Any()).Return(tt.createErr).AnyTimes()

			eventUsecase := NewEventUsecase(mockEventRepository)

			_, err := eventUsecase.CreateEvent(tt.userID, tt.title, tt.description, tt.startTime, tt.endTime)
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
		name            string
		success         bool
		userID          string
		listByUserIDErr error
	}{
		{"success list events", true, "6d322c66-bf4d-427a-970c-874f3745f653", nil},
		{"failure empty user id", false, "", nil},
		{"failure list by user id error", false, "6d322c66-bf4d-427a-970c-874f3745f653", errors.New("list by user id error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockEvent := mocks.NewMockEvent(ctrl)
			mockEventRepository := mocks.NewMockEventRepository(ctrl)
			mockEventRepository.EXPECT().ListByUserID(gomock.Any()).Return([]event.Event{mockEvent}, tt.listByUserIDErr).AnyTimes()

			eventUsecase := NewEventUsecase(mockEventRepository)

			_, err := eventUsecase.ListEvents(tt.userID)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}
