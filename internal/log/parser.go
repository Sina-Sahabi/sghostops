package log

import (
	"strings"
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
	return makeValidEntry(parts)
}
