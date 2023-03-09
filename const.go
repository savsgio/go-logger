package logger

import "time"

// Logger levels.
const (
	invalid Level = iota - 1
	PRINT
	TRACE
	FATAL
	ERROR
	WARNING
	INFO
	DEBUG
)

// Logger flags.
const (
	Ldatetime Flag = 1 << iota
	Ltimestamp
	LUTC
	Llongfile
	Lshortfile
	LstdFlags = Ldatetime
)

const calldepth = 6

const unknownFile = "???"

const (
	printLevelStr   = ""
	traceLevelStr   = "TRACE"
	fatalLevelStr   = "FATAL"
	errorLevelStr   = "ERROR"
	warningLevelStr = "WARNING"
	infoLevelStr    = "INFO"
	debugLevelStr   = "DEBUG"
)

const defaultTextSeparator = " - "

const (
	defaultJSONFieldKeyDatetime  = "datetime"
	defaultJSONFieldKeyTimestamp = "timestamp"
	defaultJSONFieldKeyLevel     = "level"
	defaultJSONFieldKeyFile      = "file"
	defaultJSONFieldKeyMessage   = "message"
)

const defaultDatetimeLayout = time.RFC3339
