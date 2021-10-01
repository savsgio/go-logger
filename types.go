package logger

import (
	"io"
	"sync"
)

type Level int

type Flag int

type Field struct {
	Key   string
	Value interface{}
}

type Config struct {
	Level  Level
	Output io.Writer
	Fields []Field

	Flag      Flag
	Datetime  bool
	Timestamp bool
	UTC       bool
	Shortfile bool
	Longfile  bool

	calldepth int
}

type Logger struct {
	mu      sync.RWMutex // ensures atomic writes; protects the following fields
	cfg     Config
	encoder Encoder
}

type Encoder interface {
	Copy() Encoder
	SetConfig(cfg Config)
	Encode(level, msg string, args []interface{}) error
}

type EncoderBase struct {
	cfg           Config
	fieldsEncoded string
}

type EncoderText struct {
	EncoderBase
}

type EncoderJSON struct {
	EncoderBase
}
