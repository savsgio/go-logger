package logger

import (
	"github.com/savsgio/gotils/strings"
)

// NewEncoderJSON creates a new json encoder.
func NewEncoderJSON(cfg EncoderJSONConfig) *EncoderJSON {
	if cfg.FieldMap.DatetimeKey == "" {
		cfg.FieldMap.DatetimeKey = defaultJSONFieldKeyDatetime
	}

	if cfg.FieldMap.TimestampKey == "" {
		cfg.FieldMap.TimestampKey = defaultJSONFieldKeyTimestamp
	}

	if cfg.FieldMap.LevelKey == "" {
		cfg.FieldMap.LevelKey = defaultJSONFieldKeyLevel
	}

	if cfg.FieldMap.FileKey == "" {
		cfg.FieldMap.FileKey = defaultJSONFieldKeyFile
	}

	if cfg.FieldMap.FunctionKey == "" {
		cfg.FieldMap.FunctionKey = defaultJSONFieldKeyFunction
	}

	if cfg.FieldMap.MessageKey == "" {
		cfg.FieldMap.MessageKey = defaultJSONFieldKeyMessage
	}

	if cfg.DatetimeLayout == "" {
		cfg.DatetimeLayout = defaultDatetimeLayout
	}

	if cfg.TimestampFormat == 0 {
		cfg.TimestampFormat = defaultTimestampFormat
	}

	enc := new(EncoderJSON)
	enc.cfg = cfg

	return enc
}

// Copy returns a copy of the json encoder.
func (enc *EncoderJSON) Copy() Encoder {
	copyEnc := NewEncoderJSON(enc.cfg)
	copyEnc.EncoderBase = *enc.EncoderBase.Copy()

	return copyEnc
}

func (enc *EncoderJSON) keys(cfg Config) (keys []string) {
	if cfg.Datetime {
		keys = append(keys, enc.cfg.FieldMap.DatetimeKey)
	}

	if cfg.Timestamp {
		keys = append(keys, enc.cfg.FieldMap.TimestampKey)
	}

	keys = append(keys, enc.cfg.FieldMap.LevelKey)

	if cfg.Shortfile || cfg.Longfile {
		keys = append(keys, enc.cfg.FieldMap.FileKey)
	}

	if cfg.Function {
		keys = append(keys, enc.cfg.FieldMap.FunctionKey)
	}

	keys = append(keys, enc.cfg.FieldMap.MessageKey)

	return keys
}

// Configure configures then encoder.
//
// - Encondes and sets the fields.
func (enc *EncoderJSON) Configure(cfg Config) {
	if len(cfg.Fields) == 0 {
		enc.SetFieldsEncoded("")

		return
	}

	buf := AcquireBuffer()
	keys := enc.keys(cfg)

	for _, field := range cfg.Fields {
		buf.WriteString("\"") // nolint:errcheck

		if strings.Include(keys, field.Key) {
			buf.WriteString("fields.") // nolint:errcheck
		}

		n := buf.Len()
		buf.WriteString(field.Key) // nolint:errcheck
		buf.Escape(n)

		buf.WriteString("\":\"") // nolint:errcheck

		n = buf.Len()
		buf.WriteInterface(field.Value)
		buf.Escape(n)

		buf.WriteString("\",") // nolint:errcheck
	}

	enc.SetFieldsEncoded(buf.String())

	ReleaseBuffer(buf)
}

// Encode encodes the given entry to the buffer.
func (enc *EncoderJSON) Encode(buf *Buffer, e Entry) error { // nolint:funlen
	buf.WriteByte('{') // nolint:errcheck

	if e.Config.Datetime {
		buf.WriteString("\"")                         // nolint:errcheck
		buf.WriteString(enc.cfg.FieldMap.DatetimeKey) // nolint:errcheck
		buf.WriteString("\":\"")                      // nolint:errcheck
		buf.WriteDatetime(e.Time, enc.cfg.DatetimeLayout)
		buf.WriteString("\",") // nolint:errcheck
	}

	if e.Config.Timestamp {
		buf.WriteString("\"")                          // nolint:errcheck
		buf.WriteString(enc.cfg.FieldMap.TimestampKey) // nolint:errcheck
		buf.WriteString("\":\"")                       // nolint:errcheck
		buf.WriteTimestamp(e.Time, enc.cfg.TimestampFormat)
		buf.WriteString("\",") // nolint:errcheck
	}

	if levelStr := e.Level.String(); levelStr != "" {
		buf.WriteString("\"")                      // nolint:errcheck
		buf.WriteString(enc.cfg.FieldMap.LevelKey) // nolint:errcheck
		buf.WriteString("\":\"")                   // nolint:errcheck
		buf.WriteString(levelStr)                  // nolint:errcheck
		buf.WriteString("\",")                     // nolint:errcheck
	}

	if e.Config.Shortfile || e.Config.Longfile {
		buf.WriteString("\"")                     // nolint:errcheck
		buf.WriteString(enc.cfg.FieldMap.FileKey) // nolint:errcheck
		buf.WriteString("\":\"")                  // nolint:errcheck
		buf.WriteFileCaller(e.Caller, e.Config.Shortfile)
		buf.WriteString("\",") // nolint:errcheck
	}

	if e.Config.Function {
		buf.WriteString("\"")                         // nolint:errcheck
		buf.WriteString(enc.cfg.FieldMap.FunctionKey) // nolint:errcheck
		buf.WriteString("\":\"")                      // nolint:errcheck
		buf.WriteString(e.Caller.Function)            // nolint:errcheck
		buf.WriteString("\",")                        // nolint:errcheck
	}

	buf.WriteString(enc.FieldsEncoded())         // nolint:errcheck
	buf.WriteString("\"")                        // nolint:errcheck
	buf.WriteString(enc.cfg.FieldMap.MessageKey) // nolint:errcheck
	buf.WriteString("\":\"")                     // nolint:errcheck

	n := buf.Len()
	buf.WriteString(e.Message) // nolint:errcheck
	buf.Escape(n)

	buf.WriteString("\"}") // nolint:errcheck
	buf.WriteNewLine()

	return nil
}
