package logger

import (
	"bytes"
	"strconv"
	"time"

	gstrconv "github.com/savsgio/gotils/strconv"
	"github.com/valyala/bytebufferpool"
)

// NewEncoderJSON creates a new json encoder.
func NewEncoderJSON() *EncoderJSON {
	return new(EncoderJSON)
}

func (enc *EncoderJSON) hasBytesSpecialChars(value []byte) bool {
	if bytes.IndexByte(value, '"') >= 0 || bytes.IndexByte(value, '\\') >= 0 {
		return true
	}

	for i := 0; i < len(value); i++ {
		if value[i] < 0x20 {
			return true
		}
	}

	return false
}

func (enc *EncoderJSON) writeEscapedBytes(buf *bytebufferpool.ByteBuffer, b []byte) {
	str := bytebufferpool.Get()
	str.Set(b) // NOTE: Use as a copy of b.

	str.B = strconv.AppendQuote(str.B, gstrconv.B2S(str.B))

	buf.Write(str.B[len(b)+1 : str.Len()-1]) // nolint:errcheck
	bytebufferpool.Put(str)
}

func (enc *EncoderJSON) escape(buf *bytebufferpool.ByteBuffer, startAt int) {
	if b := buf.B[startAt:]; enc.hasBytesSpecialChars(b) {
		buf.Set(buf.B[:startAt])
		enc.writeEscapedBytes(buf, b)
	}
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

// WriteInterface writes an interface value to the buffer.
func (enc *EncoderJSON) WriteInterface(buf *bytebufferpool.ByteBuffer, value interface{}) {
	before := buf.Len()

	enc.EncoderBase.WriteInterface(buf, value)
	enc.escape(buf, before)
}

// WriteMessage writes the given message and arguments to the buffer.
func (enc *EncoderJSON) WriteMessage(buf *bytebufferpool.ByteBuffer, msg string, args []interface{}) {
	before := buf.Len()

	enc.EncoderBase.WriteMessage(buf, msg, args)
	enc.escape(buf, before)
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
