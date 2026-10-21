package analyzer

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/Sina-Sahabi/sghostops/internal/log"
)

type Pattern struct {
	ID       string
	Name     string
	Priority int
	Match    func(entry log.Entry) bool
}

type Config struct {
	Rules []RuleConfig `json:"rules"`
}

type RuleConfig struct {
	ID       string      `json:"id"`
	Name     string      `json:"name"`
	Priority int         `json:"priority"`
	Match    MatchConfig `json:"match"`
}

type MatchConfig struct {
	Levels   []log.Severity `json:"levels,omitempty"`
	Contains string         `json:"contains,omitempty"`
}

func ExtractPatterns(path string) ([]Pattern, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse json: %w", err)
	}

	patterns := []Pattern{
		// {
		// 	ID:       "?",
		// 	Type:     "?",
		// 	Priority: ?,
		// 	Match: func(e log.Entry) bool {
		// 		return e.Level == ? &&
		// 			strings.Contains(e.Message, "?")
		// 	},
		// }, ...
	}

	for _, rc := range cfg.Rules {
		matchCfg := rc.Match

		matchFunc := func(e log.Entry) bool {
			if len(matchCfg.Levels) > 0 && !slices.Contains(matchCfg.Levels, e.Level) {
				return false
			}

			if matchCfg.Contains != "" && !strings.Contains(e.Message, matchCfg.Contains) {
				return false
			}

			return true
		}

		patterns = append(patterns, Pattern{
			ID:       rc.ID,
			Name:     rc.Name,
			Priority: rc.Priority,
			Match:    matchFunc,
		})
	}

	return patterns, nil
}
