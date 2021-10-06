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

func WithFields(fields ...Field) *Logger {
	l := std.WithFields(fields...)
	l.setCalldepth(calldepth)

	return l
}

func SetFields(fields ...Field) {
	std.SetFields(fields...)
}

func SetFlags(flag Flag) {
	std.SetFlags(flag)
}

func SetLevel(level Level) {
	std.SetLevel(level)
}

func SetOutput(output io.Writer) {
	std.SetOutput(output)
}

func SetEncoder(enc Encoder) {
	std.SetEncoder(enc)
}

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
