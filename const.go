package logger

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
