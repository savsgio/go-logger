package logger

import (
	"io"
	"runtime"
	"sync"
	"time"

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

// Buffer provides the byte buffer used by encoders to encode the output.
type Buffer struct {
	b1 bytebufferpool.ByteBuffer
	b2 bytebufferpool.ByteBuffer
}

// Entry collects all the information for the output.
type Entry struct {
	Config  Config
	Time    time.Time
	Level   Level
	Caller  runtime.Frame
	Message string
}

// Config is the logger configuration.
type Config struct {
	Fields    []Field
	Datetime  bool
	Timestamp bool
	UTC       bool
	Shortfile bool
	Longfile  bool

	flag      Flag
	calldepth int
}

// Logger type.
type Logger struct {
	mu           sync.RWMutex // ensures atomic writes; protects the following fields
	cfg          Config
	level        Level
	output       io.Writer
	encoder      Encoder
	encodeOutput encodeOutputFunc
	exit         exitFunc
}

// Encoder is the interface of encoders.
type Encoder interface {
	Copy() Encoder
	FieldsEnconded() string
	SetFieldsEnconded(fieldsEncoded string)
	SetFields(fields []Field)
	Encode(*Buffer, Entry) error
}

// EncoderBase is the base of encoders.
type EncoderBase struct {
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
