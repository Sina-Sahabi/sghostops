package report

import "github.com/Sina-Sahabi/sghostops/internal/rule"

type Report struct {
	TotalEntries     int              `json:"total_entries"`
	ValidEntries     int              `json:"valid_entries"`
	MalformedEntries int              `json:"malformed_entries"`
	CountsByLevel    map[string]int   `json:"counts_by_severity"`
	CountsByService  map[string]int   `json:"counts_by_service"`
	Incidents        []rule.Incidents `json:"incidents"`
}
