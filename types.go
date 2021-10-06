package logger

import (
	"io"
	"sync"

	"github.com/valyala/bytebufferpool"
)

type encodeOutputFunc func(level Level, levelStr, msg string, args []interface{})

type Level int

type Flag int

type Field struct {
	Key   string
	Value interface{}
}

type EncoderConfig struct {
	Fields    []Field
	Flag      Flag
	Datetime  bool
	Timestamp bool
	UTC       bool
	Shortfile bool
	Longfile  bool

	calldepth int
}

type Logger struct {
	mu           sync.RWMutex // ensures atomic writes; protects the following fields
	cfg          EncoderConfig
	level        Level
	output       io.Writer
	encoder      Encoder
	encodeOutput encodeOutputFunc
}

type Encoder interface {
	Copy() Encoder
	Config() EncoderConfig
	SetConfig(cfg EncoderConfig)
	Encode(buf *bytebufferpool.ByteBuffer, levelStr, msg string, args []interface{}) error
}

type EncoderBase struct {
	cfg           EncoderConfig
	fieldsEncoded string
}

type EncoderText struct {
	EncoderBase
}

type EncoderJSON struct {
	EncoderBase
}
