package log

import "time"

type Severity string

const (
	SeverityInfo  Severity = "INFO"
	SeverityWarn  Severity = "WARN"
	SeverityError Severity = "ERROR"
	SeverityFatal Severity = "FATAL"
)

func AllSeverities() []Severity {
	return []Severity{SeverityInfo, SeverityWarn, SeverityError, SeverityFatal}
}

type Entry struct {
	Timestamp time.Time
	Level     Severity
	Service   string
	Message   string
}
