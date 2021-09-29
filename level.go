package logger

import "strings"

func ParseLevel(levelStr string) (level Level, err error) {
	switch strings.ToUpper(levelStr) {
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
