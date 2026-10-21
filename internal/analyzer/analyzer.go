package analyzer

import (
	"bufio"
	"cmp"
	"io"
	"slices"

	"github.com/Sina-Sahabi/sghostops/internal/log"
	rpt "github.com/Sina-Sahabi/sghostops/internal/report"
	"github.com/Sina-Sahabi/sghostops/internal/rule"
)

type Analyzer struct {
	parser log.Parser
	engine rule.Engine
}

func NewAnalyzer(p log.Parser, eg rule.Engine) *Analyzer {
	return &Analyzer{
		parser: p,
		engine: eg,
	}
}

func (a *Analyzer) AnalyzeStream(r io.Reader, f func(log.Entry) bool) (rpt.Report, error) {
	report := rpt.Report{
		CountsByLevel: map[string]int{
			string(log.SeverityInfo):  0,
			string(log.SeverityWarn):  0,
			string(log.SeverityError): 0,
			string(log.SeverityFatal): 0,
		},
		CountsByService: make(map[string]int),
		Incidents:       make([]rule.Incidents, 0),
	}

	incidentMap := make(map[string]*rule.Incidents)

	scanner := bufio.NewScanner(r)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	for scanner.Scan() {
		report.TotalEntries++
		rawLine := scanner.Text()

		entry, err := a.parser.Parse(rawLine)
		if err != nil {
			report.MalformedEntries++
			continue
		}

		if pattern, matched := a.engine.Evaluate(entry); !matched && f(entry) {
			report.ValidEntries++
			report.CountsByLevel[string(entry.Level)]++
			report.CountsByService[entry.Service]++

			if group, exists := incidentMap[pattern.ID]; exists {
				group.Count++
				if entry.Timestamp.Before(group.FirstOccur) {
					group.FirstOccur = entry.Timestamp
					group.FirstError = entry.Message
				}
				if entry.Timestamp.After(group.LastOccur) {
					group.LastOccur = entry.Timestamp
				}
			} else {
				incidentMap[pattern.ID] = &rule.Incidents{
					PatternID:  pattern.ID,
					Name:       pattern.Name,
					Severity:   entry.Level,
					Service:    entry.Service,
					Count:      1,
					FirstOccur: entry.Timestamp,
					LastOccur:  entry.Timestamp,
					FirstError: entry.Message,
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return rpt.Report{}, err
	}

	for _, group := range incidentMap {
		report.Incidents = append(report.Incidents, *group)
	}

	slices.SortFunc(report.Incidents, func(a, b rule.Incidents) int {
		as, _ := a.Severity.IsValid()
		bs, _ := b.Severity.IsValid()
		if c := cmp.Compare(bs, as); c != 0 {
			return c
		}
		if c := cmp.Compare(b.Count, a.Count); c != 0 {
			return c
		}
		if c := cmp.Compare(a.PatternID, b.PatternID); c != 0 {
			return c
		}
		return cmp.Compare(a.Service, b.Service)
	})

	return report, nil
}
