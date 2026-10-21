package incident

import (
	"time"

	"github.com/Sina-Sahabi/sghostops/internal/model/log"
)

type Incidents struct {
	PatternID   string       `json:"pattern_id"`
	Type        string       `json:"type"`
	Severity    log.Severity `json:"severity"`
	Service     string       `json:"service"`
	Count       int          `json:"count"`
	FirstOccur  time.Time    `json:"first_occur"`
	LastOccur   time.Time    `json:"last_occur"`
	SampleError string       `json:"sample_message"`
}
