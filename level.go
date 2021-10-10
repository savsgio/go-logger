package logger

import "strings"

// ParseLevel returns the Level constant from the given level string.
func ParseLevel(levelStr string) (level Level, err error) {
	switch strings.ToUpper(levelStr) {
	case printLevelStr:
		level = PRINT
	case traceLevelStr:
		level = TRACE
	case fatalLevelStr:
		level = FATAL
	case errorLevelStr:
		level = ERROR
	case warningLevelStr:
		level = WARNING
	case infoLevelStr:
		level = INFO
	case debugLevelStr:
		level = DEBUG
	default:
		level = invalid
		err = ErrInvalidLevel
	}

	return level, err
}
