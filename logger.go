package logger

import (
	"io"
	"os"

	"github.com/valyala/bytebufferpool"
)

func newEncodeOutputFunc(l *Logger) encodeOutputFunc {
	return func(level Level, levelStr, msg string, args []interface{}) {
		l.mu.RLock()

		if l.isLevelEnabled(level) {
			buf := bytebufferpool.Get()

			l.encoder.Encode(buf, levelStr, msg, args) // nolint:errcheck
			l.output.Write(buf.Bytes())                // nolint:errcheck

			bytebufferpool.Put(buf)
		}

		l.mu.RUnlock()
	}
}

// New create new instance of Logger.
func New(level Level, output io.Writer, fields ...Field) *Logger {
	cfg := EncoderConfig{
		calldepth: calldepth,
	}

	enc := NewEncoderText()
	enc.SetConfig(cfg)

	l := new(Logger)
	l.cfg = cfg
	l.level = level
	l.output = output
	l.encoder = enc
	l.encodeOutput = newEncodeOutputFunc(l)
	l.exit = os.Exit

	l.SetFields(fields...)
	l.SetFlags(LstdFlags)

	return l
}

func (l *Logger) getField(key string) *Field {
	for i := range l.cfg.Fields {
		field := &l.cfg.Fields[i]

		if field.Key == key {
			return field
		}
	}

	return nil
}

func (l *Logger) setCalldepth(value int) {
	l.cfg.calldepth = value
	l.encoder.SetConfig(l.cfg)
}

func (l *Logger) setFields(fields ...Field) {
	for _, field := range fields {
		if optField := l.getField(field.Key); optField != nil {
			optField.Value = field.Value
		} else {
			l.cfg.Fields = append(l.cfg.Fields, field)
		}
	}

	l.encoder.SetConfig(l.cfg)
}

func (l *Logger) isLevelEnabled(level Level) bool {
	return l.level >= level
}

func (l *Logger) clone() *Logger {
	cfgFields := make([]Field, len(l.cfg.Fields))
	copy(cfgFields, l.cfg.Fields)

	l2 := new(Logger)
	l2.cfg = l.cfg
	l2.cfg.Fields = cfgFields
	l2.level = l.level
	l2.output = l.output
	l2.encoder = l.encoder.Copy()
	l2.encodeOutput = l.encodeOutput
	l2.exit = l.exit

	return l2
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

// SetLogFlags sets the output flags for the logger.
func (l *Logger) SetFlags(flag Flag) {
	l.mu.Lock()

	l.cfg.Flag = flag
	l.cfg.Datetime = flag&Ldatetime != 0
	l.cfg.Timestamp = flag&Ltimestamp != 0
	l.cfg.UTC = flag&LUTC != 0
	l.cfg.Shortfile = flag&Lshortfile != 0
	l.cfg.Longfile = flag&Llongfile != 0

	l.encoder.SetConfig(l.cfg)

	l.mu.Unlock()
}

// SetLevel set level of log.
func (l *Logger) SetLevel(level Level) {
	l.mu.Lock()
	l.level = level
	l.mu.Unlock()
}

// SetOutput set output of log.
func (l *Logger) SetOutput(output io.Writer) {
	l.mu.Lock()
	l.output = output
	l.mu.Unlock()
}

// SetOutput set output of log.
func (l *Logger) SetEncoder(enc Encoder) {
	l.mu.Lock()
	l.encoder = enc
	l.encoder.SetConfig(l.cfg)
	l.mu.Unlock()
}

func (l *Logger) IsLevelEnabled(level Level) bool {
	l.mu.RLock()
	enabled := l.isLevelEnabled(level)
	l.mu.RUnlock()

	return enabled
}

func (l *Logger) Print(msg ...interface{}) {
	l.encodeOutput(PRINT, printLevelStr, "", msg)
}

func (l *Logger) Printf(msg string, args ...interface{}) {
	l.encodeOutput(PRINT, printLevelStr, msg, args)
}

func (l *Logger) Fatal(msg ...interface{}) {
	l.encodeOutput(FATAL, fatalLevelStr, "", msg)
	l.exit(1)
}

func (l *Logger) Fatalf(msg string, args ...interface{}) {
	l.encodeOutput(FATAL, fatalLevelStr, msg, args)
	l.exit(1)
}

func (l *Logger) Error(msg ...interface{}) {
	l.encodeOutput(ERROR, errorLevelStr, "", msg)
}

func (l *Logger) Errorf(msg string, args ...interface{}) {
	l.encodeOutput(ERROR, errorLevelStr, msg, args)
}

func (l *Logger) Warning(msg ...interface{}) {
	l.encodeOutput(WARNING, warningLevelStr, "", msg)
}

func (l *Logger) Warningf(msg string, args ...interface{}) {
	l.encodeOutput(WARNING, warningLevelStr, msg, args)
}

func (l *Logger) Info(msg ...interface{}) {
	l.encodeOutput(INFO, infoLevelStr, "", msg)
}

func (l *Logger) Infof(msg string, args ...interface{}) {
	l.encodeOutput(INFO, infoLevelStr, msg, args)
}

func (l *Logger) Debug(msg ...interface{}) {
	l.encodeOutput(DEBUG, debugLevelStr, "", msg)
}

func (l *Logger) Debugf(msg string, args ...interface{}) {
	l.encodeOutput(DEBUG, debugLevelStr, msg, args)
}
