package event

import "errors"

var (
	ErrEventNotFound    = errors.New("event not found")
	ErrPermissionDenied = errors.New("permission denied")
)
