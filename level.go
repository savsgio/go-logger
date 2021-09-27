package logger

func ParseLevel(levelStr string) Level {
	switch levelStr {
	case "fatal":
		return FATAL
	case "error":
		return ERROR
	case "warning":
		return WARNING
	case "info":
		return INFO
	case "debug":
		return DEBUG
	default:
		return INVALID
	}
}
