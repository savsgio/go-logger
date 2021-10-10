package logger

import (
	"time"

	"github.com/valyala/bytebufferpool"
)

// NewEncoderJSON creates a new json encoder.
func NewEncoderJSON() *EncoderJSON {
	return new(EncoderJSON)
}

// Copy returns a copy of the json encoder.
func (enc *EncoderJSON) Copy() Encoder {
	copyEnc := NewEncoderJSON()
	copyEnc.EncoderBase = *enc.EncoderBase.Copy()

	return copyEnc
}

// SetConfig sets the encoder config and encode the fields.
func (enc *EncoderJSON) SetConfig(cfg EncoderConfig) {
	enc.EncoderBase.SetConfig(cfg)

	buf := bytebufferpool.Get()

	for _, field := range enc.cfg.Fields {
		buf.WriteString("\"")      // nolint:errcheck
		buf.WriteString(field.Key) // nolint:errcheck
		buf.WriteString("\":\"")   // nolint:errcheck
		enc.WriteInterface(buf, field.Value)
		buf.WriteString("\",") // nolint:errcheck
	}

	enc.SetFieldsEnconded(buf.String())

	bytebufferpool.Put(buf)
}

// Encode encodes the given level string, message and arguments to the buffer.
func (enc *EncoderJSON) Encode(buf *bytebufferpool.ByteBuffer, levelStr, msg string, args []interface{}) error {
	now := time.Now()
	if enc.cfg.UTC {
		now = now.UTC()
	}

	buf.WriteByte('{') // nolint:errcheck

	if enc.cfg.Datetime {
		buf.WriteString("\"datetime\":\"") // nolint:errcheck
		enc.WriteDatetime(buf, now)
		buf.WriteString("\",") // nolint:errcheck
	}

	if enc.cfg.Timestamp {
		buf.WriteString("\"timestamp\":\"") // nolint:errcheck
		enc.WriteTimestamp(buf, now)
		buf.WriteString("\",") // nolint:errcheck
	}

	if levelStr != "" {
		buf.WriteString("\"level\":\"") // nolint:errcheck
		buf.WriteString(levelStr)       // nolint:errcheck
		buf.WriteString("\",")          // nolint:errcheck
	}

	if enc.cfg.Shortfile || enc.cfg.Longfile {
		buf.WriteString("\"file\":\"") // nolint:errcheck
		enc.WriteFileCaller(buf)
		buf.WriteString("\",") // nolint:errcheck
	}

	enc.WriteFieldsEnconded(buf)

	buf.WriteString("\"message\":\"") // nolint:errcheck
	enc.WriteMessage(buf, msg, args)
	buf.WriteString("\"}") // nolint:errcheck

	enc.WriteNewLine(buf)

	return nil
}
