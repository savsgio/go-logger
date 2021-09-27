package logger

import (
	"io"
	"sync"
)

type Logger struct {
	mu      sync.RWMutex // ensures atomic writes; protects the following fields
	level   Level
	flag    int
	output  io.Writer
	options Options
	encoder Encoder

	fatalEnabled   bool
	errorEnabled   bool
	warningEnabled bool
	infoEnabled    bool
	debugEnabled   bool
}

type Level int

type Field struct {
	Key   string
	Value interface{}
}

type Options struct {
	Fields           []Field
	UTC              bool
	Date             bool
	Time             bool
	TimeMicroseconds bool
	Shortfile        bool
	Longfile         bool

	calldepth int
}

type Encoder interface {
	SetOutput(output io.Writer)
	SetOptions(opts Options)
	Encode(level, msg string, args []interface{}) error
}

type EncoderBase struct {
	output io.Writer
	opts   Options

	fieldsEncoded string
}

type EncoderText struct {
	EncoderBase
}

type EncoderJSON struct {
	EncoderBase
}
