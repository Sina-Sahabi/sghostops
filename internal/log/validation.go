package log

import (
	"strings"
	"time"
)

func (s Severity) IsValid() (int, bool) {
	switch s {
	case SeverityFatal:
		return 3, true
	case SeverityError:
		return 2, true
	case SeverityWarn:
		return 1, true
	case SeverityInfo:
		return 0, true
	default:
		return -1, false
	}
}

func makeValidEntry(parts []string) (Entry, error) {
	if len(parts) < 4 {
		return Entry{}, ErrMissingComponents
	}

	tsStr, levelStr, service, message := parts[0], parts[1], parts[2], strings.TrimSpace(parts[3])
	if service == "" || message == "" {
		return Entry{}, ErrMissingComponents
	}

	ts, err := time.Parse(time.RFC3339, tsStr)
	if err != nil {
		return Entry{}, ErrInvalidTimestamp
	}

	level := Severity(levelStr)
	if _, ok := level.IsValid(); !ok {
		return Entry{}, ErrInvalidSeverity
	}

	return Entry{
		Timestamp: ts.UTC(),
		Level:     level,
		Service:   service,
		Message:   message,
	}, nil
}
