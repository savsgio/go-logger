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
	Lfunction

	LstdFlags = Ldatetime
)

// Logger timestamp formats.
const (
	TimestampFormatSeconds TimestampFormat = iota + 1
	TimestampFormatNanoseconds
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
	defaultJSONFieldKeyFunction  = "func"
	defaultJSONFieldKeyMessage   = "message"
)

const defaultDatetimeLayout = time.RFC3339

const defaultTimestampFormat = TimestampFormatSeconds
