package logger

import (
	"io"
	"sync"

	"github.com/valyala/bytebufferpool"
)

type encodeOutputFunc func(level Level, msg string, args []interface{})

type exitFunc func(code int)

// Level type.
type Level int

// Flag type.
type Flag int

// Field type.
type Field struct {
	Key   string
	Value interface{}
}

// EncoderConfig is the encoder configuration.
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

// Logger type.
type Logger struct {
	mu           sync.RWMutex // ensures atomic writes; protects the following fields
	cfg          EncoderConfig
	level        Level
	output       io.Writer
	encoder      Encoder
	encodeOutput encodeOutputFunc
	exit         exitFunc
}

// Encoder is the interface of encoders.
type Encoder interface {
	Copy() Encoder
	Config() EncoderConfig
	SetConfig(cfg EncoderConfig)
	Encode(buf *bytebufferpool.ByteBuffer, levelStr, msg string, args []interface{}) error
}

// EncoderBase is the base of encoders.
type EncoderBase struct {
	cfg           EncoderConfig
	fieldsEncoded string
}

// EncoderText is the text enconder.
type EncoderText struct {
	EncoderBase
}

// EncoderJSON is the json encoder.
type EncoderJSON struct {
	EncoderBase
}
