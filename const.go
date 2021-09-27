package logger

const calldepth = 6

const (
	INVALID Level = iota - 1
	FATAL
	ERROR
	WARNING
	INFO
	DEBUG
)

const (
	fatalLevelStr   = "FATAL"
	errorLevelStr   = "ERROR"
	warningLevelStr = "WARNING"
	infoLevelStr    = "INFO"
	debugLevelStr   = "DEBUG"
)
