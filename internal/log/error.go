package log

import "errors"

var (
	ErrMalformedEntry    = errors.New("malformed log entry")
	ErrInvalidTimestamp  = errors.New("invalid timestamp format")
	ErrInvalidSeverity   = errors.New("unsupported severity level")
	ErrMissingComponents = errors.New("missing required log fields")
)
