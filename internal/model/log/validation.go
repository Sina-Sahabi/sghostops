package log

func (s Severity) IsValid() (int, bool) {
	switch s {
	case SeverityFatal:
		return 3, true
	case SeverityError:
		return 2, true
	case SeverityWarn:
		return 1, true
	case SeverityInfo:
		return 0, true
	default:
		return -1, false
	}
}
