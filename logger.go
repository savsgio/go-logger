package logger

import (
	"fmt"
	"io"
	"log"
	"os"
)

// New create new instance of Logger.
func New(name string, level string, output io.Writer) *Logger {
	enc := NewEncoderText()

	l := new(Logger)
	l.name = name

	l.SetEncoder(enc)
	l.SetLevel(level)
	l.SetFlags(log.LstdFlags)
	l.SetOutput(output)

	return l
}

func (l *Logger) encode(level, msg string, args []interface{}) {
	l.mu.RLock()
	l.encoder.Encode(level, msg, args) // nolint:errcheck
	l.mu.RUnlock()
}

func (l *Logger) checkLevel(level int) bool {
	return l.level >= level
}

// SetLevel set level of log.
func (l *Logger) SetLevel(level string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	switch level {
	case FATAL:
		l.level = fatalLevel
	case ERROR:
		l.level = errorLevel
	case WARNING:
		l.level = warningLevel
	case INFO:
		l.level = infoLevel
	case DEBUG:
		l.level = debugLevel
	default:
		panic(fmt.Sprintf("Invalid log level, only can use {%s|%s|%s|%s|%s}", FATAL, ERROR, WARNING, INFO, DEBUG))
	}

	l.fatalEnabled = l.checkLevel(fatalLevel)
	l.errorEnabled = l.checkLevel(errorLevel)
	l.warningEnabled = l.checkLevel(warningLevel)
	l.infoEnabled = l.checkLevel(infoLevel)
	l.debugEnabled = l.checkLevel(debugLevel)

	l.encoder.SetOptions(l.options)
}

// SetLogFlags sets the output flags for the logger.
func (l *Logger) SetFlags(flag int) {
	l.mu.Lock()

	l.options.UTC = flag&log.LUTC != 0
	l.options.Date = flag&log.Ldate != 0
	l.options.Time = flag&log.Ltime != 0
	l.options.TimeMicroseconds = flag&log.Lmicroseconds != 0
	l.options.Shortfile = flag&log.Lshortfile != 0
	l.options.Longfile = flag&log.Llongfile != 0

	l.encoder.SetOptions(l.options)

	l.mu.Unlock()
}

// SetOutput set output of log.
func (l *Logger) SetOutput(output io.Writer) {
	l.mu.Lock()
	l.output = output
	l.encoder.SetOutput(output)
	l.mu.Unlock()
}

// SetOutput set output of log.
func (l *Logger) SetEncoder(enc Encoder) {
	l.mu.Lock()
	l.encoder = enc
	l.encoder.SetOutput(l.output)
	l.encoder.SetOptions(l.options)
	l.mu.Unlock()
}

func (l *Logger) FatalEnabled() bool {
	return l.fatalEnabled
}

func (l *Logger) Fatal(msg ...interface{}) {
	l.encode(fatalPrefix, "", msg)
	os.Exit(1)
}

func (l *Logger) Fatalf(msg string, args ...interface{}) {
	l.encode(fatalPrefix, msg, args)
	os.Exit(1)
}

func (l *Logger) ErrorEnabled() bool {
	return l.errorEnabled
}

func (l *Logger) Error(msg ...interface{}) {
	if l.ErrorEnabled() {
		l.encode(errorPrefix, "", msg)
	}
}

func (l *Logger) Errorf(msg string, args ...interface{}) {
	if l.ErrorEnabled() {
		l.encode(errorPrefix, msg, args)
	}
}

func (l *Logger) WarningEnabled() bool {
	return l.warningEnabled
}

func (l *Logger) Warning(msg ...interface{}) {
	if l.WarningEnabled() {
		l.encode(warningPrefix, "", msg)
	}
}

func (l *Logger) Warningf(msg string, args ...interface{}) {
	if l.WarningEnabled() {
		l.encode(warningPrefix, msg, args)
	}
}

func (l *Logger) InfoEnabled() bool {
	return l.infoEnabled
}

func (l *Logger) Info(msg ...interface{}) {
	if l.InfoEnabled() {
		l.encode(infoPrefix, "", msg)
	}
}

func (l *Logger) Infof(msg string, args ...interface{}) {
	if l.InfoEnabled() {
		l.encode(infoPrefix, msg, args)
	}
}

func (l *Logger) DebugEnabled() bool {
	return l.debugEnabled
}

func (l *Logger) Debug(msg ...interface{}) {
	if l.DebugEnabled() {
		l.encode(debugPrefix, "", msg)
	}
}

func (l *Logger) Debugf(msg string, args ...interface{}) {
	if l.DebugEnabled() {
		l.encode(debugPrefix, msg, args)
	}
}

func (l *Logger) Print(msg ...interface{}) {
	l.encode("", "", msg)
}

func (l *Logger) Printf(msg string, args ...interface{}) {
	l.encode("", msg, args)
}
