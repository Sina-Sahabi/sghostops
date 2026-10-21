package analyzer

type Engine struct {
	patterns []Pattern
}

func NewEngine(patterns []Pattern) *Engine {
	return &Engine{patterns: patterns}
}
