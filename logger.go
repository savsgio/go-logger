package logger

import (
	"io"
	"log"
	"os"
)

func newLogger(level Level, output io.Writer, enc Encoder, flag int, field ...Field) *Logger {
	l := new(Logger)
	l.SetEncoder(enc)
	l.SetLevel(level)
	l.SetFlags(flag)
	l.SetOutput(output)
	l.SetFields(field...)
	l.setCalldepth(calldepth)

	return l
}

// New create new instance of Logger.
func New(name string, level Level, output io.Writer) *Logger {
	enc := NewEncoderText()

	fields := make([]Field, 0)
	if name != "" {
		fields = append(fields, Field{"name", name})
	}

	return newLogger(level, output, enc, log.LstdFlags, fields...)
}

func (l *Logger) checkLevel(level Level) bool {
	return l.level >= level
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

func (l *Logger) encode(level, msg string, args []interface{}) {
	l.mu.RLock()
	l.encoder.Encode(level, msg, args) // nolint:errcheck
	l.mu.RUnlock()
}

func (l *Logger) WithFields(fields ...Field) *Logger {
	l.mu.RLock()

	l2 := newLogger(l.level, l.output, l.encoder, l.flag, l.options.Fields...)
	l2.SetFields(fields...)

	l.mu.RUnlock()

	return l2
}

func (l *Logger) SetFields(fields ...Field) {
	l.mu.Lock()

	for _, field := range fields {
		if optField := l.getField(field.Key); optField != nil {
			optField.Value = field.Value
		} else {
			l.options.Fields = append(l.options.Fields, field)
		}
	}

	l.encoder.SetOptions(l.options)

	l.mu.Unlock()
}

// SetLevel set level of log.
func (l *Logger) SetLevel(level Level) {
	l.mu.Lock()

	l.level = level
	l.fatalEnabled = l.checkLevel(FATAL)
	l.errorEnabled = l.checkLevel(ERROR)
	l.warningEnabled = l.checkLevel(WARNING)
	l.infoEnabled = l.checkLevel(INFO)
	l.debugEnabled = l.checkLevel(DEBUG)

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

func (l *Logger) FatalEnabled() bool {
	return l.fatalEnabled
}

func (l *Logger) Fatal(msg ...interface{}) {
	l.encode(fatalLevelStr, "", msg)
	os.Exit(1)
}

func (l *Logger) Fatalf(msg string, args ...interface{}) {
	l.encode(fatalLevelStr, msg, args)
	os.Exit(1)
}

func (l *Logger) ErrorEnabled() bool {
	return l.errorEnabled
}

func (l *Logger) Error(msg ...interface{}) {
	if l.ErrorEnabled() {
		l.encode(errorLevelStr, "", msg)
	}
}

func (l *Logger) Errorf(msg string, args ...interface{}) {
	if l.ErrorEnabled() {
		l.encode(errorLevelStr, msg, args)
	}
}

func (l *Logger) WarningEnabled() bool {
	return l.warningEnabled
}

func (l *Logger) Warning(msg ...interface{}) {
	if l.WarningEnabled() {
		l.encode(warningLevelStr, "", msg)
	}
}

func (l *Logger) Warningf(msg string, args ...interface{}) {
	if l.WarningEnabled() {
		l.encode(warningLevelStr, msg, args)
	}
}

func (l *Logger) InfoEnabled() bool {
	return l.infoEnabled
}

func (l *Logger) Info(msg ...interface{}) {
	if l.InfoEnabled() {
		l.encode(infoLevelStr, "", msg)
	}
}

func (l *Logger) Infof(msg string, args ...interface{}) {
	if l.InfoEnabled() {
		l.encode(infoLevelStr, msg, args)
	}
}

func (l *Logger) DebugEnabled() bool {
	return l.debugEnabled
}

func (l *Logger) Debug(msg ...interface{}) {
	if l.DebugEnabled() {
		l.encode(debugLevelStr, "", msg)
	}
}

func (l *Logger) Debugf(msg string, args ...interface{}) {
	if l.DebugEnabled() {
		l.encode(debugLevelStr, msg, args)
	}
}

func (l *Logger) Print(msg ...interface{}) {
	l.encode("", "", msg)
}

func (l *Logger) Printf(msg string, args ...interface{}) {
	l.encode("", msg, args)
}
