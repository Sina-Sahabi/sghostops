package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/Sina-Sahabi/sghostops/internal/analyzer"
	"github.com/Sina-Sahabi/sghostops/internal/log"
	"github.com/Sina-Sahabi/sghostops/internal/report"
	"github.com/Sina-Sahabi/sghostops/internal/rule"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	fs := flag.NewFlagSet("ghostops", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	var jsonOutPath, fileOutPath, logFilePath, timeLowerBound, timeUpperBound, severity, service, cfg string

	fs.StringVar(&jsonOutPath, "o", "", "optional - path to export report as JSON")
	fs.StringVar(&fileOutPath, "w", "", "optional - path to export report as text file")
	fs.StringVar(&timeLowerBound, "tlb", "", "optional - report the incidents after the specified time")
	fs.StringVar(&timeUpperBound, "tub", "", "optional - report the incidents before the specified time")
	fs.StringVar(&severity, "", "svt", "optional - report the incidents with the specified severity")
	fs.StringVar(&service, "", "svc", "optional - report the incidents for the specified service")
	fs.StringVar(&cfg, "", "cfg", "optional - report the incidents for the specified set of rules")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] <input-file>\n\nFlags:\n", os.Args[0])
		flag.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return err
	}

	var tfo, tel time.Time
	if timeLowerBound != "" {
		var err error
		tfo, err = time.Parse(time.RFC3339, timeLowerBound)
		if err != nil {
			fs.Usage()
			return errors.New("lower bound time should be RFC3339 format")
		}
	}
	if timeUpperBound != "" {
		var err error
		tel, err = time.Parse(time.RFC3339, timeUpperBound)
		if err != nil {
			fs.Usage()
			return errors.New("upper bound time should be RFC3339 format")
		}
	}
	var sev log.Severity
	if severity != "" {
		sev = log.Severity(severity)
		if _, ok := sev.IsValid(); !ok {
			fs.Usage()
			return fmt.Errorf("Invalid Severity - Choose %#v", log.AllSeverities())
		}
	}

	if fs.NArg() != 1 {
		fs.Usage()
		return errors.New("expected exactly one log file path argument")
	}

	logFilePath = fs.Arg(0)

	inFile, err := os.Open(logFilePath)
	if err != nil {
		return fmt.Errorf("failed to open log file %q: %w", logFilePath, err)
	}
	defer inFile.Close()

	patterns, err := rule.ExtractPatterns(cfg)
	if err != nil {
		return fmt.Errorf("failed to load rules file %q: %w", logFilePath, err)
	}
	TextParser := log.NewTextParser()
	Engine := rule.NewEngine(patterns)
	streamAnalyzer := analyzer.NewAnalyzer(TextParser, *Engine)

	limitFunc := func(e log.Entry) bool {
		if timeLowerBound != "" && e.Timestamp.Before(tfo) {
			return false
		}
		if timeUpperBound != "" && e.Timestamp.After(tel) {
			return false
		}
		if severity != "" && e.Level != sev {
			return false
		}
		if service != "" && e.Service != service {
			return false
		}
		return true
	}

	rep, err := streamAnalyzer.AnalyzeStream(inFile, limitFunc)
	if err != nil {
		return fmt.Errorf("failed while analyzing log stream: %w", err)
	}

	termReporter := report.NewTerminalReporter()
	humanOut := os.Stdout
	if fileOutPath != "" {
		var err error
		humanOut, err = os.Create(logFilePath)
		if err != nil {
			return fmt.Errorf("failed to open human-readable output file %q: %w", fileOutPath, err)
		}
		defer inFile.Close()
	}
	if err := termReporter.Write(humanOut, rep); err != nil {
		return fmt.Errorf("failed to write human-readable report: %w", err)
	}

	if jsonOutPath != "" {
		if err := writeJSONFile(jsonOutPath, rep); err != nil {
			return fmt.Errorf("failed to export JSON report to %q: %w", jsonOutPath, err)
		}
	}

	return nil
}

func writeJSONFile(path string, rep report.Report) (retErr error) {
	outFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := outFile.Close(); closeErr != nil && retErr == nil {
			retErr = closeErr
		}
	}()

	jsonReporter := report.NewJSONReporter()
	return jsonReporter.Write(outFile, rep)
}
