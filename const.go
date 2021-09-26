package logger

const calldepth = 6

const (
	fatalLevel = iota
	errorLevel
	warningLevel
	infoLevel
	debugLevel
)

const (
	fatalPrefix   = "FATAL"
	errorPrefix   = "ERROR"
	warningPrefix = "WARNING"
	infoPrefix    = "INFO"
	debugPrefix   = "DEBUG"
)

const (
	FATAL   = "fatal"
	ERROR   = "error"
	WARNING = "warning"
	INFO    = "info"
	DEBUG   = "debug"
)
