package logger

import "time"

// Logger levels.
const (
	invalid Level = iota - 1
	PRINT
	PANIC
	FATAL
	ERROR
	WARNING
	INFO
	DEBUG
	TRACE
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

const (
	calldepth    = 4
	calldepthStd = calldepth + 1
)

const unknownFile = "???"

const (
	printLevelStr   = ""
	panicLevelStr   = "PANIC"
	fatalLevelStr   = "FATAL"
	errorLevelStr   = "ERROR"
	warningLevelStr = "WARNING"
	infoLevelStr    = "INFO"
	debugLevelStr   = "DEBUG"
	traceLevelStr   = "TRACE"
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
