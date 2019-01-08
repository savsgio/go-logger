package logger

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/valyala/bytebufferpool"
)

// New create new instance of Logger
func New(name string, level string, output io.Writer) *Logger {
	l := &Logger{name: name, out: output}
	l.instance = log.New(output, "", log.LstdFlags)

	l.SetLevel(level)

	return l
}

func (l *Logger) isStd() bool {
	return l.name == stdName
}

// Check level to make print or not
func (l *Logger) checkLevel(level int) bool {
	return l.level >= level
}

// Get complete prefix if name of the logger isn't 'std'
func (l *Logger) writePrefix(buff msgBuffer, prefix string) {
	buff.SetString("- ")

	if !l.isStd() {
		buff.WriteString(l.name)
		buff.WriteString(" - ")
	}

	if prefix != "" {
		buff.WriteString(prefix)
		buff.WriteString(" - ")
	}
}

func (l *Logger) output(prefix string, msg ...interface{}) {
	buff := bytebufferpool.Get()
	defer bytebufferpool.Put(buff)

	l.writePrefix(buff, prefix)
	fmt.Fprint(buff, msg...)

	l.instance.Output(calldepth, buff.String())
}

func (l *Logger) outputf(prefix string, msg string, v ...interface{}) {
	buff := bytebufferpool.Get()
	defer bytebufferpool.Put(buff)

	l.writePrefix(buff, prefix)
	fmt.Fprintf(buff, msg, v...)

	l.instance.Output(calldepth, buff.String())
}

// SetLevel set level of log
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
}

// SetOutput set output of log
func (l *Logger) SetOutput(output io.Writer) {
	l.mu.Lock()
	l.out = output
	l.mu.Unlock()

	l.instance.SetOutput(output)
}

func (l *Logger) FatalEnabled() bool {
	return l.fatalEnabled
}

func (l *Logger) Fatal(msg ...interface{}) {
	l.output(fatalPrefix, msg...)
	os.Exit(1)
}

func (l *Logger) Fatalf(msg string, v ...interface{}) {
	l.outputf(fatalPrefix, msg, v...)
	os.Exit(1)
}

func (l *Logger) ErrorEnabled() bool {
	return l.errorEnabled
}

func (l *Logger) Error(msg ...interface{}) {
	if l.ErrorEnabled() {
		l.output(errorPrefix, msg...)
	}
}

func (l *Logger) Errorf(msg string, v ...interface{}) {
	if l.ErrorEnabled() {
		l.outputf(errorPrefix, msg, v...)
	}
}

func (l *Logger) WarningEnabled() bool {
	return l.warningEnabled
}

func (l *Logger) Warning(msg ...interface{}) {
	if l.WarningEnabled() {
		l.output(warningPrefix, msg...)
	}
}

func (l *Logger) Warningf(msg string, v ...interface{}) {
	if l.WarningEnabled() {
		l.outputf(warningPrefix, msg, v...)
	}
}

func (l *Logger) InfoEnabled() bool {
	return l.infoEnabled
}

func (l *Logger) Info(msg ...interface{}) {
	if l.InfoEnabled() {
		l.output(infoPrefix, msg...)
	}
}

func (l *Logger) Infof(msg string, v ...interface{}) {
	if l.InfoEnabled() {
		l.outputf(infoPrefix, msg, v...)
	}
}

func (l *Logger) DebugEnabled() bool {
	return l.debugEnabled
}

func (l *Logger) Debug(msg ...interface{}) {
	if l.DebugEnabled() {
		l.output(debugPrefix, msg...)
	}
}

func (l *Logger) Debugf(msg string, v ...interface{}) {
	if l.DebugEnabled() {
		l.outputf(debugPrefix, msg, v...)
	}
}

func (l *Logger) Print(msg ...interface{}) {
	l.output("", msg...)
}

func (l *Logger) Printf(msg string, v ...interface{}) {
	l.outputf("", msg, v...)
}
