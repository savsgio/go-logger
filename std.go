package logger

import (
	"io"
	"os"
)

var std = newStd()

func newStd() *Logger {
	l := New(INFO, os.Stderr)
	l.setCalldepth(calldepth + 1)

	return l
}

// WithFields returns a copy of the standard logger with the given fields.
func WithFields(fields ...Field) *Logger {
	l := std.WithFields(fields...)
	l.setCalldepth(calldepth)

	return l
}

// SetFields sets the fields to the standard logger.
func SetFields(fields ...Field) {
	std.SetFields(fields...)
}

// SetFlags sets the output flags to the standard logger.
func SetFlags(flag Flag) {
	std.SetFlags(flag)
}

// SetLevel sets the level to the standard logger.
func SetLevel(level Level) {
	std.SetLevel(level)
}

// SetOutput sets the output to the standard logger.
func SetOutput(output io.Writer) {
	std.SetOutput(output)
}

// SetEncoder sets the encoder to the standard logger.
func SetEncoder(enc Encoder) {
	std.SetEncoder(enc)
}

// IsLevelEnabled checks if the given level is enabled on the standard logger.
func IsLevelEnabled(level Level) bool {
	return std.IsLevelEnabled(level)
}

func Print(msg ...interface{}) {
	std.Print(msg...)
}

func Printf(msg string, args ...interface{}) {
	std.Printf(msg, args...)
}

func Fatal(msg ...interface{}) {
	std.Fatal(msg...)
}

func Fatalf(msg string, args ...interface{}) {
	std.Fatalf(msg, args...)
}

func Error(msg ...interface{}) {
	std.Error(msg...)
}

func Errorf(msg string, args ...interface{}) {
	std.Errorf(msg, args...)
}

func Warning(msg ...interface{}) {
	std.Warning(msg...)
}

func Warningf(msg string, args ...interface{}) {
	std.Warningf(msg, args...)
}

func Info(msg ...interface{}) {
	std.Info(msg...)
}

func Infof(msg string, args ...interface{}) {
	std.Infof(msg, args...)
}

func Debug(msg ...interface{}) {
	std.Debug(msg...)
}

func Debugf(msg string, args ...interface{}) {
	std.Debugf(msg, args...)
}
