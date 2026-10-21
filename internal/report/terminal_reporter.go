package report

import (
	"cmp"
	"fmt"
	"io"
	"slices"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/Sina-Sahabi/sghostops/internal/log"
)

type TerminalReporter struct{}

func NewTerminalReporter() *TerminalReporter {
	return &TerminalReporter{}
}

func (tr *TerminalReporter) Write(w io.Writer, rep Report) error {
	fmt.Fprintf(w, "%s\n", strings.Repeat("=", 82))
	fmt.Fprintf(w, "%sGHOSTOPS LOG ANALYSIS REPORT%s\n", strings.Repeat("=", 26), strings.Repeat("=", 26))
	fmt.Fprintf(w, "%s\n", strings.Repeat("=", 82))

	fmt.Fprintf(w, "1. OVERVIEW\n")
	fmt.Fprintf(w, "-----------\n")
	fmt.Fprintf(w, "  Total Entries Processed : %d\n", rep.TotalEntries)
	fmt.Fprintf(w, "  Valid Log Entries       : %d\n", rep.ValidEntries)
	fmt.Fprintf(w, "  Malformed Entries       : %d\n\n", rep.MalformedEntries)

	fmt.Fprintf(w, "2. ENTRIES BY SEVERITY\n")
	fmt.Fprintf(w, "----------------------\n")
	for _, sev := range log.AllSeverities() {
		count := rep.CountsByLevel[string(sev)]
		fmt.Fprintf(w, "  %-7s : %d\n", sev, count)
	}
	fmt.Fprintf(w, "\n")

	fmt.Fprintf(w, "3. ENTRIES BY SERVICE\n")
	fmt.Fprintf(w, "---------------------\n")
	if len(rep.CountsByService) == 0 {
		fmt.Fprintf(w, "  No valid service entries found.\n\n")
	} else {
		type serviceCount struct {
			service string
			count   int
		}
		services := make([]serviceCount, 0, len(rep.CountsByService))
		for svc, count := range rep.CountsByService {
			services = append(services, serviceCount{service: svc, count: count})
		}
		slices.SortFunc(services, func(a, b serviceCount) int {
			return cmp.Compare(a.service, b.service)
		})
		for _, sc := range services {
			fmt.Fprintf(w, "  %-22s : %d\n", sc.service, sc.count)
		}
		fmt.Fprintf(w, "\n")
	}

	fmt.Fprintf(w, "4. DETECTED INCIDENT GROUPS (%d)\n", len(rep.Incidents))
	fmt.Fprintf(w, "--------------------------------\n")
	if len(rep.Incidents) == 0 {
		fmt.Fprintf(w, "  No incidents detected.\n")
	} else {
		tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
		_, _ = fmt.Fprintln(tw, "  RULE\tSEVERITY\tSERVICE\tCOUNT\tFIRST OCCUR(UTC)\tLAST OCCUR(UTC)\tINCIDENT TYPE")
		_, _ = fmt.Fprintln(tw, "  ----\t--------\t-------\t-----\t----------------\t---------------\t-------------")
		for _, inc := range rep.Incidents {
			_, _ = fmt.Fprintf(
				tw,
				"  %s\t%s\t%s\t%d\t%s\t%s\t%s\n",
				inc.PatternID,
				inc.Severity,
				inc.Service,
				inc.Count,
				inc.FirstOccur.Format(time.RFC3339),
				inc.LastOccur.Format(time.RFC3339),
				inc.Name,
			)
		}
		if flushErr := tw.Flush(); flushErr != nil {
			return flushErr
		}
	}

	fmt.Fprintf(w, "\n%s\n", strings.Repeat("=", 80))
	return nil
}
