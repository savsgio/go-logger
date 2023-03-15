package logger

import (
	"strings"
)

// ParseLevel returns the Level constant from the given level string.
func ParseLevel(levelStr string) (level Level, err error) {
	switch strings.ToUpper(levelStr) {
	case printLevelStr:
		level = PRINT
	case traceLevelStr:
		level = TRACE
	case panicLevelStr:
		level = PANIC
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

// Strings returns the string representation of the level.
func (l Level) String() string {
	switch l {
	case PRINT:
		return printLevelStr
	case TRACE:
		return traceLevelStr
	case PANIC:
		return panicLevelStr
	case FATAL:
		return fatalLevelStr
	case ERROR:
		return errorLevelStr
	case WARNING:
		return warningLevelStr
	case INFO:
		return infoLevelStr
	case DEBUG:
		return debugLevelStr
	case invalid:
		fallthrough
	default:
		return ErrInvalidLevel.Error()
	}
}
