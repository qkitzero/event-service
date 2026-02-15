package event

import "errors"

var (
	ErrEventNotFound     = errors.New("event not found")
	ErrPermissionDenied  = errors.New("permission denied")
	ErrStartTimeRequired = errors.New("start time is required")
	ErrEndTimeRequired   = errors.New("end time is required")
)
