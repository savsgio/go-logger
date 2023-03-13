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

type levelHooks struct {
	store     map[Level][]Hook
	errOutput io.Writer
}

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

// TimestampFormat type.
type TimestampFormat int

// Entry collects all the information for the output.
type Entry struct {
	Config     Config
	Time       time.Time
	Level      Level
	Caller     runtime.Frame
	Message    string
	RawMessage string
	Args       []interface{}
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
	hooks        *levelHooks
	exit         exitFunc
}

// Hook represents a extended functionality that will be fired when logging.
//
// NOTE: This is not run concurrently, so be quite with locks.
type Hook interface {
	// Levels returns the levels at which the hook fires.
	Levels() []Level

	// Fire is hook function.
	//
	// NOTE: The returned error will be written to `os.Stderr`.
	Fire(Entry) error
}

// Encoder represents the encoders contract.
type Encoder interface {
	Copy() Encoder
	FieldsEncoded() string
	SetFieldsEncoded(fieldsEncoded string)
	Configure(cfg Config)
	Encode(*Buffer, Entry) error
}

// EncoderBase is the base of encoders.
type EncoderBase struct {
	fieldsEncoded string
}

// EncoderTextConfig is the configuration of text encoder.
type EncoderTextConfig struct {
	// Default: -
	Separator string

	// Default: time.RFC3339
	DatetimeLayout string

	// Default: TimestampFormatSeconds
	TimestampFormat TimestampFormat
}

// EncoderText is the text enconder.
type EncoderText struct {
	EncoderBase

	cfg EncoderTextConfig
}

// EncoderJSON is the json encoder.
type EncoderJSON struct {
	EncoderBase

	cfg EncoderJSONConfig
}

// EncoderJSONConfig is the configuration of json encoder.
type EncoderJSONConfig struct {
	FieldMap EnconderJSONFieldMap

	// Default: time.RFC3339
	DatetimeLayout string

	// Default: TimestampFormatSeconds
	TimestampFormat TimestampFormat
}

// EnconderJSONFieldMap defines name of keys.
type EnconderJSONFieldMap struct {
	// Default: datetime
	DatetimeKey string

	// Default: timestamp
	TimestampKey string

	// Default: level
	LevelKey string

	// Default: file
	FileKey string

	// Default: message
	MessageKey string
}
