package logger

const calldepth = 6

const (
	invalid Level = iota - 1
	PRINT
	FATAL
	ERROR
	WARNING
	INFO
	DEBUG
)

const (
	printLevelStr   = ""
	fatalLevelStr   = "FATAL"
	errorLevelStr   = "ERROR"
	warningLevelStr = "WARNING"
	infoLevelStr    = "INFO"
	debugLevelStr   = "DEBUG"
)
