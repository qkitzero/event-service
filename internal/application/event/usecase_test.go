package event

import (
	"errors"
	"testing"

	"go.uber.org/mock/gomock"

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
		createErr   error
	}{
		{"success create event", true, "6d322c66-bf4d-427a-970c-874f3745f653", "title", "description", nil},
		{"failure empty user id", false, "", "title", "description", nil},
		{"failure empty title", false, "6d322c66-bf4d-427a-970c-874f3745f653", "", "description", nil},
		{"failure empty description", false, "6d322c66-bf4d-427a-970c-874f3745f653", "title", "", nil},
		{"failure create error", false, "6d322c66-bf4d-427a-970c-874f3745f653", "title", "description", errors.New("create error")},
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

			_, err := eventUsecase.CreateEvent(tt.userID, tt.title, tt.description)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}
