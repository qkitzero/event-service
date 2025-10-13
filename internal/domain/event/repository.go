package event

import "github.com/qkitzero/event-service/internal/domain/user"

type EventRepository interface {
	Create(event Event) error
	Update(event Event) error
	FindByID(id EventID) (Event, error)
	ListByUserID(userID user.UserID) ([]Event, error)
	Delete(id EventID) error
}
