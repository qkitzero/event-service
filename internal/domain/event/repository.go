package event

import (
	"context"

	"github.com/qkitzero/event-service/internal/domain/user"
)

type EventRepository interface {
	Create(ctx context.Context, event Event) error
	Update(ctx context.Context, event Event) error
	FindByID(ctx context.Context, id EventID) (Event, error)
	FindAllByUserID(ctx context.Context, userID user.UserID) ([]Event, error)
	Delete(ctx context.Context, id EventID) error
}
