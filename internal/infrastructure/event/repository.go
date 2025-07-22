package event

import (
	"gorm.io/gorm"

	"github.com/qkitzero/event-service/internal/domain/event"
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
			Title:       e.Title(),
			Description: e.Description(),
			CreatedAt:   e.CreatedAt(),
		}

		if err := tx.Create(&eventModel).Error; err != nil {
			return err
		}

		return nil
	})
}
