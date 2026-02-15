package event

import (
	"context"
	"errors"
	"testing"
	"time"

	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/qkitzero/event-service/internal/domain/event"
	domainuser "github.com/qkitzero/event-service/internal/domain/user"
	mocksauth "github.com/qkitzero/event-service/mocks/application/auth"
	mocksuser "github.com/qkitzero/event-service/mocks/application/user"
	mocks "github.com/qkitzero/event-service/mocks/domain/event"
)

func TestCreateEvent(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		success     bool
		ctx         context.Context
		userID      string
		getUserErr  error
		title       string
		description string
		startTime   *timestamppb.Timestamp
		endTime     *timestamppb.Timestamp
		color       string
		createErr   error
	}{
		{"success create event", true, context.Background(), "6d322c66-bf4d-427a-970c-874f3745f653", nil, "title", "description", timestamppb.Now(), timestamppb.Now(), "#FFFFFF", nil},
		{"failure get user error", false, context.Background(), "", errors.New("get user error"), "title", "description", timestamppb.Now(), timestamppb.Now(), "#FFFFFF", nil},
		{"failure empty user id", false, context.Background(), "", nil, "title", "description", timestamppb.Now(), timestamppb.Now(), "#FFFFFF", nil},
		{"failure empty title", false, context.Background(), "6d322c66-bf4d-427a-970c-874f3745f653", nil, "", "description", timestamppb.Now(), timestamppb.Now(), "#FFFFFF", nil},
		{"failure empty description", false, context.Background(), "6d322c66-bf4d-427a-970c-874f3745f653", nil, "title", "", timestamppb.Now(), timestamppb.Now(), "#FFFFFF", nil},
		{"failure nil start time", false, context.Background(), "6d322c66-bf4d-427a-970c-874f3745f653", nil, "title", "description", nil, timestamppb.Now(), "#FFFFFF", nil},
		{"failure nil end time", false, context.Background(), "6d322c66-bf4d-427a-970c-874f3745f653", nil, "title", "description", timestamppb.Now(), nil, "#FFFFFF", nil},
		{"failure invalid color", false, context.Background(), "6d322c66-bf4d-427a-970c-874f3745f653", nil, "title", "description", timestamppb.Now(), timestamppb.Now(), "red", nil},
		{"failure create error", false, context.Background(), "6d322c66-bf4d-427a-970c-874f3745f653", nil, "title", "description", timestamppb.Now(), timestamppb.Now(), "#FFFFFF", errors.New("create error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthService := mocksauth.NewMockAuthService(ctrl)
			mockUserService := mocksuser.NewMockUserService(ctrl)
			mockEventRepository := mocks.NewMockEventRepository(ctrl)
			mockUserService.EXPECT().GetUser(tt.ctx).Return(tt.userID, tt.getUserErr).AnyTimes()
			mockEventRepository.EXPECT().Create(gomock.Any()).Return(tt.createErr).AnyTimes()

			eventUsecase := NewEventUsecase(mockAuthService, mockUserService, mockEventRepository)

			_, err := eventUsecase.CreateEvent(tt.ctx, tt.title, tt.description, tt.startTime, tt.endTime, tt.color)
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
		eventOwnerID string
		userID       string
		getUserErr   error
		eventID      string
		title        string
		description  string
		startTime    *timestamppb.Timestamp
		endTime      *timestamppb.Timestamp
		color        string
		findByIDErr  error
		updateErr    error
	}{
		{"success update event", true, context.Background(), "6d322c66-bf4d-427a-970c-874f3745f653", "6d322c66-bf4d-427a-970c-874f3745f653", nil, "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", "title", "description", timestamppb.Now(), timestamppb.Now(), "#FFFFFF", nil, nil},
		{"success update event with nil times", true, context.Background(), "6d322c66-bf4d-427a-970c-874f3745f653", "6d322c66-bf4d-427a-970c-874f3745f653", nil, "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", "title", "description", nil, nil, "#FFFFFF", nil, nil},
		{"failure get user error", false, context.Background(), "6d322c66-bf4d-427a-970c-874f3745f653", "6d322c66-bf4d-427a-970c-874f3745f653", errors.New("get user error"), "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", "title", "description", timestamppb.Now(), timestamppb.Now(), "#FFFFFF", nil, nil},
		{"failure permission denied", false, context.Background(), "6d322c66-bf4d-427a-970c-874f3745f653", "00000000-0000-0000-0000-000000000001", nil, "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", "title", "description", timestamppb.Now(), timestamppb.Now(), "#FFFFFF", nil, nil},
		{"failure empty event id", false, context.Background(), "6d322c66-bf4d-427a-970c-874f3745f653", "6d322c66-bf4d-427a-970c-874f3745f653", nil, "", "title", "description", timestamppb.Now(), timestamppb.Now(), "#FFFFFF", nil, nil},
		{"failure empty title", false, context.Background(), "6d322c66-bf4d-427a-970c-874f3745f653", "6d322c66-bf4d-427a-970c-874f3745f653", nil, "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", "", "description", timestamppb.Now(), timestamppb.Now(), "#FFFFFF", nil, nil},
		{"failure empty description", false, context.Background(), "6d322c66-bf4d-427a-970c-874f3745f653", "6d322c66-bf4d-427a-970c-874f3745f653", nil, "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", "title", "", timestamppb.Now(), timestamppb.Now(), "#FFFFFF", nil, nil},
		{"failure invalid color", false, context.Background(), "6d322c66-bf4d-427a-970c-874f3745f653", "6d322c66-bf4d-427a-970c-874f3745f653", nil, "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", "title", "description", timestamppb.Now(), timestamppb.Now(), "red", nil, nil},
		{"failure find by id error", false, context.Background(), "6d322c66-bf4d-427a-970c-874f3745f653", "6d322c66-bf4d-427a-970c-874f3745f653", nil, "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", "title", "description", timestamppb.Now(), timestamppb.Now(), "#FFFFFF", errors.New("find by id error"), nil},
		{"failure update error", false, context.Background(), "6d322c66-bf4d-427a-970c-874f3745f653", "6d322c66-bf4d-427a-970c-874f3745f653", nil, "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", "title", "description", timestamppb.Now(), timestamppb.Now(), "#FFFFFF", nil, errors.New("update error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthService := mocksauth.NewMockAuthService(ctrl)
			mockUserService := mocksuser.NewMockUserService(ctrl)
			mockUserService.EXPECT().GetUser(tt.ctx).Return(tt.userID, tt.getUserErr).AnyTimes()
			eventUserID, _ := domainuser.NewUserIDFromString(tt.eventOwnerID)
			mockEvent := mocks.NewMockEvent(ctrl)
			mockEvent.EXPECT().UserID().Return(eventUserID).AnyTimes()
			mockEvent.EXPECT().StartTime().Return(time.Now()).AnyTimes()
			mockEvent.EXPECT().EndTime().Return(time.Now()).AnyTimes()
			mockEvent.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return().AnyTimes()
			mockEventRepository := mocks.NewMockEventRepository(ctrl)
			mockEventRepository.EXPECT().FindByID(gomock.Any()).Return(mockEvent, tt.findByIDErr).AnyTimes()
			mockEventRepository.EXPECT().Update(gomock.Any()).Return(tt.updateErr).AnyTimes()

			eventUsecase := NewEventUsecase(mockAuthService, mockUserService, mockEventRepository)

			_, err := eventUsecase.UpdateEvent(tt.ctx, tt.eventID, tt.title, tt.description, tt.startTime, tt.endTime, tt.color)
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
		name         string
		success      bool
		ctx          context.Context
		eventOwnerID string
		userID       string
		getUserErr   error
		eventID      string
		findByIDErr  error
	}{
		{"success get event", true, context.Background(), "6d322c66-bf4d-427a-970c-874f3745f653", "6d322c66-bf4d-427a-970c-874f3745f653", nil, "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", nil},
		{"failure get user error", false, context.Background(), "6d322c66-bf4d-427a-970c-874f3745f653", "6d322c66-bf4d-427a-970c-874f3745f653", errors.New("get user error"), "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", nil},
		{"failure permission denied", false, context.Background(), "6d322c66-bf4d-427a-970c-874f3745f653", "00000000-0000-0000-0000-000000000001", nil, "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", nil},
		{"failure empty event id", false, context.Background(), "6d322c66-bf4d-427a-970c-874f3745f653", "6d322c66-bf4d-427a-970c-874f3745f653", nil, "", nil},
		{"failure find by id error", false, context.Background(), "6d322c66-bf4d-427a-970c-874f3745f653", "6d322c66-bf4d-427a-970c-874f3745f653", nil, "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", errors.New("find by id error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthService := mocksauth.NewMockAuthService(ctrl)
			mockUserService := mocksuser.NewMockUserService(ctrl)
			mockUserService.EXPECT().GetUser(tt.ctx).Return(tt.userID, tt.getUserErr).AnyTimes()
			eventUserID, _ := domainuser.NewUserIDFromString(tt.eventOwnerID)
			mockEvent := mocks.NewMockEvent(ctrl)
			mockEvent.EXPECT().UserID().Return(eventUserID).AnyTimes()
			mockEventRepository := mocks.NewMockEventRepository(ctrl)
			mockEventRepository.EXPECT().FindByID(gomock.Any()).Return(mockEvent, tt.findByIDErr).AnyTimes()

			eventUsecase := NewEventUsecase(mockAuthService, mockUserService, mockEventRepository)

			_, err := eventUsecase.GetEvent(tt.ctx, tt.eventID)
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

			mockAuthService := mocksauth.NewMockAuthService(ctrl)
			mockUserService := mocksuser.NewMockUserService(ctrl)
			mockUserService.EXPECT().GetUser(tt.ctx).Return(tt.userID, tt.getUserErr).AnyTimes()
			mockEvent := mocks.NewMockEvent(ctrl)
			mockEventRepository := mocks.NewMockEventRepository(ctrl)
			mockEventRepository.EXPECT().FindAllByUserID(gomock.Any()).Return([]event.Event{mockEvent}, tt.findAllByUserIDErr).AnyTimes()

			eventUsecase := NewEventUsecase(mockAuthService, mockUserService, mockEventRepository)

			_, err := eventUsecase.ListEvents(tt.ctx)
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
		name         string
		success      bool
		ctx          context.Context
		eventOwnerID string
		userID       string
		getUserErr   error
		eventID      string
		findByIDErr  error
		deleteErr    error
	}{
		{"success delete event", true, context.Background(), "6d322c66-bf4d-427a-970c-874f3745f653", "6d322c66-bf4d-427a-970c-874f3745f653", nil, "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", nil, nil},
		{"failure get user error", false, context.Background(), "6d322c66-bf4d-427a-970c-874f3745f653", "6d322c66-bf4d-427a-970c-874f3745f653", errors.New("get user error"), "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", nil, nil},
		{"failure permission denied", false, context.Background(), "6d322c66-bf4d-427a-970c-874f3745f653", "00000000-0000-0000-0000-000000000001", nil, "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", nil, nil},
		{"failure empty event id", false, context.Background(), "6d322c66-bf4d-427a-970c-874f3745f653", "6d322c66-bf4d-427a-970c-874f3745f653", nil, "", nil, nil},
		{"failure find by id error", false, context.Background(), "6d322c66-bf4d-427a-970c-874f3745f653", "6d322c66-bf4d-427a-970c-874f3745f653", nil, "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", errors.New("find by id error"), nil},
		{"failure delete error", false, context.Background(), "6d322c66-bf4d-427a-970c-874f3745f653", "6d322c66-bf4d-427a-970c-874f3745f653", nil, "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", nil, errors.New("delete error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthService := mocksauth.NewMockAuthService(ctrl)
			mockUserService := mocksuser.NewMockUserService(ctrl)
			mockUserService.EXPECT().GetUser(tt.ctx).Return(tt.userID, tt.getUserErr).AnyTimes()
			eventUserID, _ := domainuser.NewUserIDFromString(tt.eventOwnerID)
			mockEvent := mocks.NewMockEvent(ctrl)
			mockEvent.EXPECT().UserID().Return(eventUserID).AnyTimes()
			mockEventRepository := mocks.NewMockEventRepository(ctrl)
			mockEventRepository.EXPECT().FindByID(gomock.Any()).Return(mockEvent, tt.findByIDErr).AnyTimes()
			mockEventRepository.EXPECT().Delete(gomock.Any()).Return(tt.deleteErr).AnyTimes()

			eventUsecase := NewEventUsecase(mockAuthService, mockUserService, mockEventRepository)

			err := eventUsecase.DeleteEvent(tt.ctx, tt.eventID)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}
