package logger

import (
	"io"
	"os"
)

var std = New(stdName, INFO, os.Stderr)

func SetLevel(level string) {
	std.SetLevel(level)
}

func SetOutput(output io.Writer) {
	std.SetOutput(output)
}

func FatalEnabled() bool {
	return std.FatalEnabled()
}

func Fatal(msg ...interface{}) {
	std.Fatal(msg...)
}

func Fatalf(msg string, v ...interface{}) {
	std.Fatalf(msg, v...)
}

func ErrorEnabled() bool {
	return std.ErrorEnabled()
}

func Error(msg ...interface{}) {
	std.Error(msg...)
}

func Errorf(msg string, v ...interface{}) {
	std.Errorf(msg, v...)
}

func WarningEnabled() bool {
	return std.WarningEnabled()
}

func Warning(msg ...interface{}) {
	std.Warning(msg...)
}

func Warningf(msg string, v ...interface{}) {
	std.Warningf(msg, v...)
}

func InfoEnabled() bool {
	return std.InfoEnabled()
}

func Info(msg ...interface{}) {
	std.Info(msg...)
}

func Infof(msg string, v ...interface{}) {
	std.Infof(msg, v...)
}

func DebugEnabled() bool {
	return std.DebugEnabled()
}

func Debug(msg ...interface{}) {
	std.Debug(msg...)
}

func Debugf(msg string, v ...interface{}) {
	std.Debugf(msg, v...)
}

func Print(msg ...interface{}) {
	std.Print(msg...)
}

func Printf(msg string, v ...interface{}) {
	std.Printf(msg, v...)
}
