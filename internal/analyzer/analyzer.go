package analyzer

import (
	"io"

	"github.com/Sina-Sahabi/sghostops/internal/log"
	"github.com/Sina-Sahabi/sghostops/internal/report"
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

func (a *Analyzer) AnalyzeStream(r io.Reader) (report.Report, error) {
	// report := report.Report{
	// 	CountsByLevel: map[string]int{
	// 		string(log.SeverityInfo):  0,
	// 		string(log.SeverityWarn):  0,
	// 		string(log.SeverityError): 0,
	// 		string(log.SeverityFatal): 0,
	// 	},
	// 	CountsByService: make(map[string]int),
	// 	Incidents:       make([]log.Incidents, 0),
	// }
}
