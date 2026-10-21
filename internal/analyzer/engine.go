package analyzer

import "github.com/Sina-Sahabi/sghostops/internal/log"

type Engine struct {
	patterns []Pattern
}

func NewEngine(patterns []Pattern) *Engine {
	return &Engine{patterns: patterns}
}

func (e *Engine) Evaluate(entry log.Entry) (Pattern, bool) {
	for _, rule := range e.patterns {
		if rule.Match(entry) {
			return rule, true
		}
	}
	return Pattern{}, false
}
