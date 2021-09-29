package logger

import (
	"io"
	"log"
	"os"
)

func newLogger(level Level, output io.Writer, enc Encoder, flag int, fields ...Field) *Logger {
	l := new(Logger)
	l.SetEncoder(enc)
	l.SetLevel(level)
	l.SetFlags(flag)
	l.SetOutput(output)
	l.SetFields(fields...)
	l.setCalldepth(calldepth)

	return l
}

// New create new instance of Logger.
func New(level Level, output io.Writer, fields ...Field) *Logger {
	enc := NewEncoderText()

	return newLogger(INFO, os.Stderr, enc, log.LstdFlags, fields...)
}

func (l *Logger) getField(key string) *Field {
	for i := range l.options.Fields {
		field := &l.options.Fields[i]

		if field.Key == key {
			return field
		}
	}

	return nil
}

func (l *Logger) setCalldepth(calldepth int) {
	l.options.calldepth = calldepth
	l.encoder.SetOptions(l.options)
}

func (l *Logger) setFields(fields ...Field) {
	for _, field := range fields {
		if optField := l.getField(field.Key); optField != nil {
			optField.Value = field.Value
		} else {
			l.options.Fields = append(l.options.Fields, field)
		}
	}

	l.encoder.SetOptions(l.options)
}

func (l *Logger) isLevelEnabled(level Level) bool {
	return l.level >= level
}

func (l *Logger) encode(level Level, levelStr, msg string, args []interface{}) {
	l.mu.RLock()

	if l.isLevelEnabled(level) {
		l.encoder.Encode(levelStr, msg, args) // nolint:errcheck
	}

	l.mu.RUnlock()
}

func (l *Logger) clone() *Logger {
	return newLogger(l.level, l.output, l.encoder, l.flag, l.options.Fields...)
}

func (l *Logger) WithFields(fields ...Field) *Logger {
	l.mu.RLock()

	l2 := l.clone()
	l2.setFields(fields...)

	l.mu.RUnlock()

	return l2
}

func (l *Logger) SetFields(fields ...Field) {
	l.mu.Lock()
	l.setFields(fields...)
	l.mu.Unlock()
}

func (l *Logger) IsLevelEnabled(level Level) bool {
	l.mu.RLock()
	enabled := l.isLevelEnabled(level)
	l.mu.RUnlock()

	return enabled
}

// SetLevel set level of log.
func (l *Logger) SetLevel(level Level) {
	l.mu.Lock()
	l.level = level
	l.mu.Unlock()
}

// SetLogFlags sets the output flags for the logger.
func (l *Logger) SetFlags(flag int) {
	l.mu.Lock()

	l.flag = flag
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

func (l *Logger) Fatal(msg ...interface{}) {
	l.encode(FATAL, fatalLevelStr, "", msg)
	os.Exit(1)
}

func (l *Logger) Fatalf(msg string, args ...interface{}) {
	l.encode(FATAL, fatalLevelStr, msg, args)
	os.Exit(1)
}

func (l *Logger) Error(msg ...interface{}) {
	l.encode(ERROR, errorLevelStr, "", msg)
}

func (l *Logger) Errorf(msg string, args ...interface{}) {
	l.encode(ERROR, errorLevelStr, msg, args)
}

func (l *Logger) Warning(msg ...interface{}) {
	l.encode(WARNING, warningLevelStr, "", msg)
}

func (l *Logger) Warningf(msg string, args ...interface{}) {
	l.encode(WARNING, warningLevelStr, msg, args)
}

func (l *Logger) Info(msg ...interface{}) {
	l.encode(INFO, infoLevelStr, "", msg)
}

func (l *Logger) Infof(msg string, args ...interface{}) {
	l.encode(INFO, infoLevelStr, msg, args)
}

func (l *Logger) Debug(msg ...interface{}) {
	l.encode(DEBUG, debugLevelStr, "", msg)
}

func (l *Logger) Debugf(msg string, args ...interface{}) {
	l.encode(DEBUG, debugLevelStr, msg, args)
}

func (l *Logger) Print(msg ...interface{}) {
	l.encode(PRINT, printLevelStr, "", msg)
}

func (l *Logger) Printf(msg string, args ...interface{}) {
	l.encode(PRINT, printLevelStr, msg, args)
}
