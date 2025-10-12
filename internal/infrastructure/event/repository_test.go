package event

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/qkitzero/event-service/internal/domain/event"
	"github.com/qkitzero/event-service/internal/domain/user"
	mocksevent "github.com/qkitzero/event-service/mocks/domain/event"
	"github.com/qkitzero/event-service/testutil"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		setup   func(mock sqlmock.Sqlmock, event event.Event)
	}{
		{
			name:    "success create event",
			success: true,
			setup: func(mock sqlmock.Sqlmock, event event.Event) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "events" ("id","user_id","title","description","start_time","end_time","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`)).
					WithArgs(event.ID(), event.UserID(), event.Title(), event.Description(), testutil.AnyTime{}, testutil.AnyTime{}, testutil.AnyTime{}, testutil.AnyTime{}).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		{
			name:    "failure create event error",
			success: false,
			setup: func(mock sqlmock.Sqlmock, event event.Event) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "events" ("id","user_id","title","description","start_time","end_time","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`)).
					WithArgs(event.ID(), event.UserID(), event.Title(), event.Description(), testutil.AnyTime{}, testutil.AnyTime{}, testutil.AnyTime{}, testutil.AnyTime{}).
					WillReturnError(errors.New("create event error"))

				mock.ExpectRollback()
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sqlDB, mock, err := sqlmock.New()
			if err != nil {
				t.Errorf("failed to new sqlmock: %s", err)
			}

			gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
			if err != nil {
				t.Errorf("failed to open gorm: %s", err)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockEvent := mocksevent.NewMockEvent(ctrl)
			mockEvent.EXPECT().ID().Return(event.EventID{UUID: uuid.New()}).AnyTimes()
			mockEvent.EXPECT().UserID().Return(user.UserID{UUID: uuid.New()}).AnyTimes()
			mockEvent.EXPECT().Title().Return(event.Title("title")).AnyTimes()
			mockEvent.EXPECT().Description().Return(event.Description("description")).AnyTimes()
			mockEvent.EXPECT().StartTime().Return(time.Now()).AnyTimes()
			mockEvent.EXPECT().EndTime().Return(time.Now()).AnyTimes()
			mockEvent.EXPECT().CreatedAt().Return(time.Now()).AnyTimes()
			mockEvent.EXPECT().UpdatedAt().Return(time.Now()).AnyTimes()

			tt.setup(mock, mockEvent)

			repo := NewEventRepository(gormDB)

			err = repo.Create(mockEvent)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestListByUserID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		userID  user.UserID
		setup   func(mock sqlmock.Sqlmock, userID user.UserID)
	}{
		{
			name:    "success list by user id",
			success: true,
			userID:  user.UserID{UUID: uuid.New()},
			setup: func(mock sqlmock.Sqlmock, userID user.UserID) {
				eventRows := sqlmock.NewRows([]string{"id", "user_id", "title", "description", "start_time", "end_time", "created_at", "updated_at"}).
					AddRow(uuid.New(), userID, "title", "description", time.Now(), time.Now(), time.Now(), time.Now())
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "events" WHERE user_id = $1`)).
					WithArgs(userID).
					WillReturnRows(eventRows)
			},
		},
		{
			name:    "failure find events error",
			success: false,
			userID:  user.UserID{UUID: uuid.New()},
			setup: func(mock sqlmock.Sqlmock, userID user.UserID) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "events" WHERE user_id = $1`)).
					WithArgs(userID).
					WillReturnError(errors.New("find events error"))
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sqlDB, mock, err := sqlmock.New()
			if err != nil {
				t.Errorf("failed to new sqlmock: %s", err)
			}

			gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
			if err != nil {
				t.Errorf("failed to open gorm: %s", err)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tt.setup(mock, tt.userID)

			repo := NewEventRepository(gormDB)

			_, err = repo.ListByUserID(tt.userID)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
