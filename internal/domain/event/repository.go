package event

type EventRepository interface {
	Create(event Event) error
}
