package analyzer

import "github.com/Sina-Sahabi/sghostops/internal/log"

type Pattern struct {
	ID       string
	Type     string
	Priority int
	Match    func(entry log.Entry) bool
}
