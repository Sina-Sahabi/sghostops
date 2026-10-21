package log

import (
	"strings"
	"time"
)

type Parser interface {
	Parse(entry any) (Entry, error)
}

type TextParser struct{}

func NewTextParser() *TextParser {
	return &TextParser{}
}

func (p *TextParser) Parse(rawline any) (Entry, error) {
	line, ok := rawline.(string)
	if !ok {
		return Entry{}, ErrMalformedEntry
	}
	line = strings.TrimSpace(line)
	if line == "" {
		return Entry{}, ErrMalformedEntry
	}

	parts := strings.SplitN(line, " ", 4)
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
