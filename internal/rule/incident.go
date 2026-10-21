package rule

import (
	"time"

	"github.com/Sina-Sahabi/sghostops/internal/log"
)

type Incidents struct {
	PatternID  string       `json:"pattern_id"`
	Name       string       `json:"name"`
	Severity   log.Severity `json:"severity"`
	Service    string       `json:"service"`
	Count      int          `json:"count"`
	FirstOccur time.Time    `json:"first_occur"`
	LastOccur  time.Time    `json:"last_occur"`
	FirstError string       `json:"first_message"`
}
