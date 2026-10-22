Compile: go build -o bin/sghostops cmd/sghostops/main.go

Usage: ./bin/sghostops [flags] <input-file>

Flags:
  -cfg string
        optional - report the incidents for the specified set of rules
  -o string
        optional - path to export report as JSON
  -svc string
        optional - report the incidents for the specified service
  -svt string
        optional - report the incidents with the specified severity
  -tlb string
        optional - report the incidents after the specified time
  -tub string
        optional - report the incidents before the specified time

Example: ./bin/sghostops -cfg rules.json -tlb 2026-10-02T12:20:47Z -tub 2026-10-02T12:20:47Z -o res.json tests/sample.log