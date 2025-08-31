package event

import (
	"gorm.io/gorm"

	"github.com/qkitzero/event-service/internal/domain/event"
	"github.com/qkitzero/event-service/internal/domain/user"
)

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) event.EventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) Create(e event.Event) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		eventModel := EventModel{
			ID:          e.ID(),
			UserID:      e.UserID(),
			Title:       e.Title(),
			Description: e.Description(),
			StartTime:   e.StartTime(),
			EndTime:     e.EndTime(),
			CreatedAt:   e.CreatedAt(),
		}

		if err := tx.Create(&eventModel).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *eventRepository) ListByUserID(userID user.UserID) ([]event.Event, error) {
	var eventModels []EventModel
	if err := r.db.Where("user_id = ?", userID).Find(&eventModels).Error; err != nil {
		return nil, err
	}

	var events []event.Event
	for _, eventModel := range eventModels {
		e := event.NewEvent(
			eventModel.ID,
			eventModel.UserID,
			eventModel.Title,
			eventModel.Description,
			eventModel.StartTime,
			eventModel.EndTime,
			eventModel.CreatedAt,
		)
		events = append(events, e)
	}

	return events, nil
}
