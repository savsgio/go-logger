package logger

import (
	"io"
	"os"
	"time"
)

func newEncodeOutputFunc(l *Logger) encodeOutputFunc {
	return func(level Level, msg string, args []interface{}) {
		l.mu.RLock()

		if l.isLevelEnabled(level) {
			buf := AcquireBuffer()

			e := Entry{
				Config:  l.cfg,
				Level:   level,
				Message: buf.formatMessage(msg, args),
			}
			e.Caller.File = unknownFile
			e.Caller.Line = 0

			if l.cfg.Datetime || l.cfg.Timestamp {
				e.Time = time.Now()

				if l.cfg.UTC {
					e.Time = e.Time.UTC()
				}
			}

			if l.cfg.Shortfile || l.cfg.Longfile {
				e.Caller = getFileCaller(l.cfg.calldepth)
			}

			l.encoder.Encode(buf, e)    // nolint:errcheck
			l.output.Write(buf.Bytes()) // nolint:errcheck

			ReleaseBuffer(buf)
		}

		l.mu.RUnlock()
	}
}

// New creates a new Logger.
func New(level Level, output io.Writer, fields ...Field) *Logger {
	l := new(Logger)
	l.cfg = Config{
		calldepth: calldepth,
	}
	l.level = level
	l.output = output
	l.encoder = NewEncoderText()
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
}

func (l *Logger) setFields(fields ...Field) {
	for _, field := range fields {
		if optField := l.getField(field.Key); optField != nil {
			optField.Value = field.Value
		} else {
			l.cfg.Fields = append(l.cfg.Fields, field)
		}
	}

	l.encoder.SetFields(l.cfg.Fields)
}

func (l *Logger) isLevelEnabled(level Level) bool {
	return l.level >= level
}

func (l *Logger) clone() *Logger {
	l2 := new(Logger)
	l2.cfg = l.cfg.Copy()
	l2.level = l.level
	l2.output = l.output
	l2.encoder = l.encoder.Copy()
	l2.encodeOutput = newEncodeOutputFunc(l2)
	l2.exit = l.exit

	return l2
}

// WithFields returns a logger copy with the given fields.
func (l *Logger) WithFields(fields ...Field) *Logger {
	l.mu.RLock()

	l2 := l.clone()
	l2.setFields(fields...)

	l.mu.RUnlock()

	return l2
}

// SetFields sets the logger fields.
func (l *Logger) SetFields(fields ...Field) {
	l.mu.Lock()
	l.setFields(fields...)
	l.mu.Unlock()
}

// SetFlags sets the logger output flags.
func (l *Logger) SetFlags(flag Flag) {
	l.mu.Lock()

	l.cfg.Datetime = flag&Ldatetime != 0
	l.cfg.Timestamp = flag&Ltimestamp != 0
	l.cfg.UTC = flag&LUTC != 0
	l.cfg.Shortfile = flag&Lshortfile != 0
	l.cfg.Longfile = flag&Llongfile != 0
	l.cfg.flag = flag

	l.mu.Unlock()
}

// SetLevel sets the logger level.
func (l *Logger) SetLevel(level Level) {
	l.mu.Lock()
	l.level = level
	l.mu.Unlock()
}

// SetOutput sets the logger output.
func (l *Logger) SetOutput(output io.Writer) {
	l.mu.Lock()
	l.output = output
	l.mu.Unlock()
}

// SetEncoder sets the logger encoder.
func (l *Logger) SetEncoder(enc Encoder) {
	l.mu.Lock()
	l.encoder = enc
	l.encoder.SetFields(l.cfg.Fields)
	l.mu.Unlock()
}

// IsLevelEnabled checks if the given level is enabled on the logger.
func (l *Logger) IsLevelEnabled(level Level) bool {
	l.mu.RLock()
	enabled := l.isLevelEnabled(level)
	l.mu.RUnlock()

	return enabled
}

func (l *Logger) Print(msg ...interface{}) {
	l.encodeOutput(PRINT, "", msg)
}

func (l *Logger) Printf(msg string, args ...interface{}) {
	l.encodeOutput(PRINT, msg, args)
}

func (l *Logger) Trace(msg ...interface{}) {
	l.encodeOutput(TRACE, "", msg)
}

func (l *Logger) Tracef(msg string, args ...interface{}) {
	l.encodeOutput(TRACE, msg, args)
}

func (l *Logger) Fatal(msg ...interface{}) {
	l.encodeOutput(FATAL, "", msg)
	l.exit(1)
}

func (l *Logger) Fatalf(msg string, args ...interface{}) {
	l.encodeOutput(FATAL, msg, args)
	l.exit(1)
}

func (l *Logger) Error(msg ...interface{}) {
	l.encodeOutput(ERROR, "", msg)
}

func (l *Logger) Errorf(msg string, args ...interface{}) {
	l.encodeOutput(ERROR, msg, args)
}

func (l *Logger) Warning(msg ...interface{}) {
	l.encodeOutput(WARNING, "", msg)
}

func (l *Logger) Warningf(msg string, args ...interface{}) {
	l.encodeOutput(WARNING, msg, args)
}

func (l *Logger) Info(msg ...interface{}) {
	l.encodeOutput(INFO, "", msg)
}

func (l *Logger) Infof(msg string, args ...interface{}) {
	l.encodeOutput(INFO, msg, args)
}

func (l *Logger) Debug(msg ...interface{}) {
	l.encodeOutput(DEBUG, "", msg)
}

func (l *Logger) Debugf(msg string, args ...interface{}) {
	l.encodeOutput(DEBUG, msg, args)
}
