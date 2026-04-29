package event

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"

	"github.com/qkitzero/event-service/internal/domain/event"
	domainuser "github.com/qkitzero/event-service/internal/domain/user"
	mocksuser "github.com/qkitzero/event-service/mocks/application/user"
	mocksevent "github.com/qkitzero/event-service/mocks/domain/event"
)

var (
	sampleTitle, _       = event.NewTitle("title")
	sampleDescription, _ = event.NewDescription("description")
	sampleColor, _       = event.NewColor("#FFFFFF")
)

func TestCreateEvent(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name       string
		success    bool
		ctx        context.Context
		userID     string
		getUserErr error
		createErr  error
	}{
		{"success create event", true, context.Background(), "6d322c66-bf4d-427a-970c-874f3745f653", nil, nil},
		{"failure get user error", false, context.Background(), "", errors.New("get user error"), nil},
		{"failure empty user id", false, context.Background(), "", nil, nil},
		{"failure create error", false, context.Background(), "6d322c66-bf4d-427a-970c-874f3745f653", nil, errors.New("create error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserService := mocksuser.NewMockUserService(ctrl)
			mockEventRepository := mocksevent.NewMockEventRepository(ctrl)
			mockUserService.EXPECT().GetUser(tt.ctx).Return(tt.userID, tt.getUserErr).AnyTimes()
			mockEventRepository.EXPECT().Create(gomock.Any(), gomock.Any()).Return(tt.createErr).AnyTimes()

			u := NewEventUsecase(mockUserService, mockEventRepository)

			_, err := u.CreateEvent(tt.ctx, sampleTitle, sampleDescription, now, now, sampleColor)
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
	eventID, _ := event.NewEventIDFromString("fe8c2263-bbac-4bb9-a41d-b04f5afc4425")
	now := time.Now()
	tests := []struct {
		name        string
		success     bool
		ctx         context.Context
		eventUserID string
		userID      string
		getUserErr  error
		startTime   *time.Time
		endTime     *time.Time
		findByIDErr error
		updateErr   error
	}{
		{"success update event", true, context.Background(), "6d322c66-bf4d-427a-970c-874f3745f653", "6d322c66-bf4d-427a-970c-874f3745f653", nil, &now, &now, nil, nil},
		{"success update event with nil times", true, context.Background(), "6d322c66-bf4d-427a-970c-874f3745f653", "6d322c66-bf4d-427a-970c-874f3745f653", nil, nil, nil, nil, nil},
		{"failure get user error", false, context.Background(), "6d322c66-bf4d-427a-970c-874f3745f653", "6d322c66-bf4d-427a-970c-874f3745f653", errors.New("get user error"), &now, &now, nil, nil},
		{"failure permission denied", false, context.Background(), "6d322c66-bf4d-427a-970c-874f3745f653", "00000000-0000-0000-0000-000000000001", nil, &now, &now, nil, nil},
		{"failure find by id error", false, context.Background(), "6d322c66-bf4d-427a-970c-874f3745f653", "6d322c66-bf4d-427a-970c-874f3745f653", nil, &now, &now, errors.New("find by id error"), nil},
		{"failure update error", false, context.Background(), "6d322c66-bf4d-427a-970c-874f3745f653", "6d322c66-bf4d-427a-970c-874f3745f653", nil, &now, &now, nil, errors.New("update error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserService := mocksuser.NewMockUserService(ctrl)
			mockUserService.EXPECT().GetUser(tt.ctx).Return(tt.userID, tt.getUserErr).AnyTimes()
			mockEvent := mocksevent.NewMockEvent(ctrl)
			mockEvent.EXPECT().UserID().Return(domainuser.UserID{UUID: uuid.MustParse(tt.eventUserID)}).AnyTimes()
			mockEvent.EXPECT().StartTime().Return(now).AnyTimes()
			mockEvent.EXPECT().EndTime().Return(now).AnyTimes()
			mockEvent.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return().AnyTimes()
			mockEventRepository := mocksevent.NewMockEventRepository(ctrl)
			mockEventRepository.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(mockEvent, tt.findByIDErr).AnyTimes()
			mockEventRepository.EXPECT().Update(gomock.Any(), gomock.Any()).Return(tt.updateErr).AnyTimes()

			u := NewEventUsecase(mockUserService, mockEventRepository)

			_, err := u.UpdateEvent(tt.ctx, eventID, sampleTitle, sampleDescription, tt.startTime, tt.endTime, sampleColor)
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
	eventID, _ := event.NewEventIDFromString("fe8c2263-bbac-4bb9-a41d-b04f5afc4425")
	tests := []struct {
		name        string
		success     bool
		ctx         context.Context
		eventUserID string
		userID      string
		getUserErr  error
		findByIDErr error
	}{
		{"success get event", true, context.Background(), "6d322c66-bf4d-427a-970c-874f3745f653", "6d322c66-bf4d-427a-970c-874f3745f653", nil, nil},
		{"failure get user error", false, context.Background(), "6d322c66-bf4d-427a-970c-874f3745f653", "6d322c66-bf4d-427a-970c-874f3745f653", errors.New("get user error"), nil},
		{"failure permission denied", false, context.Background(), "6d322c66-bf4d-427a-970c-874f3745f653", "00000000-0000-0000-0000-000000000001", nil, nil},
		{"failure find by id error", false, context.Background(), "6d322c66-bf4d-427a-970c-874f3745f653", "6d322c66-bf4d-427a-970c-874f3745f653", nil, errors.New("find by id error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserService := mocksuser.NewMockUserService(ctrl)
			mockUserService.EXPECT().GetUser(tt.ctx).Return(tt.userID, tt.getUserErr).AnyTimes()
			mockEvent := mocksevent.NewMockEvent(ctrl)
			mockEvent.EXPECT().UserID().Return(domainuser.UserID{UUID: uuid.MustParse(tt.eventUserID)}).AnyTimes()
			mockEventRepository := mocksevent.NewMockEventRepository(ctrl)
			mockEventRepository.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(mockEvent, tt.findByIDErr).AnyTimes()

			u := NewEventUsecase(mockUserService, mockEventRepository)

			_, err := u.GetEvent(tt.ctx, eventID)
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
		name               string
		success            bool
		ctx                context.Context
		userID             string
		getUserErr         error
		findAllByUserIDErr error
	}{
		{"success list events", true, context.Background(), "6d322c66-bf4d-427a-970c-874f3745f653", nil, nil},
		{"failure get user error", false, context.Background(), "", errors.New("get user error"), nil},
		{"failure empty user id", false, context.Background(), "", nil, nil},
		{"failure find all by user id error", false, context.Background(), "6d322c66-bf4d-427a-970c-874f3745f653", nil, errors.New("find all by user id error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserService := mocksuser.NewMockUserService(ctrl)
			mockUserService.EXPECT().GetUser(tt.ctx).Return(tt.userID, tt.getUserErr).AnyTimes()
			mockEvent := mocksevent.NewMockEvent(ctrl)
			mockEventRepository := mocksevent.NewMockEventRepository(ctrl)
			mockEventRepository.EXPECT().FindAllByUserID(gomock.Any(), gomock.Any()).Return([]event.Event{mockEvent}, tt.findAllByUserIDErr).AnyTimes()

			u := NewEventUsecase(mockUserService, mockEventRepository)

			_, err := u.ListEvents(tt.ctx)
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
	eventID, _ := event.NewEventIDFromString("fe8c2263-bbac-4bb9-a41d-b04f5afc4425")
	tests := []struct {
		name        string
		success     bool
		ctx         context.Context
		eventUserID string
		userID      string
		getUserErr  error
		findByIDErr error
		deleteErr   error
	}{
		{"success delete event", true, context.Background(), "6d322c66-bf4d-427a-970c-874f3745f653", "6d322c66-bf4d-427a-970c-874f3745f653", nil, nil, nil},
		{"failure get user error", false, context.Background(), "6d322c66-bf4d-427a-970c-874f3745f653", "6d322c66-bf4d-427a-970c-874f3745f653", errors.New("get user error"), nil, nil},
		{"failure permission denied", false, context.Background(), "6d322c66-bf4d-427a-970c-874f3745f653", "00000000-0000-0000-0000-000000000001", nil, nil, nil},
		{"failure find by id error", false, context.Background(), "6d322c66-bf4d-427a-970c-874f3745f653", "6d322c66-bf4d-427a-970c-874f3745f653", nil, errors.New("find by id error"), nil},
		{"failure delete error", false, context.Background(), "6d322c66-bf4d-427a-970c-874f3745f653", "6d322c66-bf4d-427a-970c-874f3745f653", nil, nil, errors.New("delete error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserService := mocksuser.NewMockUserService(ctrl)
			mockUserService.EXPECT().GetUser(tt.ctx).Return(tt.userID, tt.getUserErr).AnyTimes()
			mockEvent := mocksevent.NewMockEvent(ctrl)
			mockEvent.EXPECT().UserID().Return(domainuser.UserID{UUID: uuid.MustParse(tt.eventUserID)}).AnyTimes()
			mockEventRepository := mocksevent.NewMockEventRepository(ctrl)
			mockEventRepository.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(mockEvent, tt.findByIDErr).AnyTimes()
			mockEventRepository.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(tt.deleteErr).AnyTimes()

			u := NewEventUsecase(mockUserService, mockEventRepository)

			err := u.DeleteEvent(tt.ctx, eventID)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}
