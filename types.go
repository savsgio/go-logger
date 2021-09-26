package logger

import (
	"io"
	"sync"
)

type Logger struct {
	mu      sync.RWMutex // ensures atomic writes; protects the following fields
	name    string
	level   int
	output  io.Writer
	options Options

	encoder Encoder
}

type Options struct {
	UTC              bool
	Date             bool
	Time             bool
	TimeMicroseconds bool
	Shortfile        bool
	Longfile         bool

	fatalEnabled   bool
	errorEnabled   bool
	warningEnabled bool
	infoEnabled    bool
	debugEnabled   bool
}

type Encoder interface {
	SetOutput(output io.Writer)
	SetOptions(opts Options)
	Encode(level, msg string, args []interface{}) error
}

type EncoderBase struct {
	output io.Writer
	opts   Options
}

type EncoderText struct {
	EncoderBase
}

type EncoderJSON struct {
	EncoderBase
}
