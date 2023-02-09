package logger

import (
	"time"

	"github.com/valyala/bytebufferpool"
)

const sepText = " - "

// NewEncoderText creates a new text encoder.
func NewEncoderText() *EncoderText {
	return new(EncoderText)
}

// Copy returns a copy of the text encoder.
func (enc *EncoderText) Copy() Encoder {
	copyEnc := NewEncoderText()
	copyEnc.EncoderBase = *enc.EncoderBase.Copy()

	return copyEnc
}

// SetConfig sets the encoder config and encode the fields.
func (enc *EncoderText) SetConfig(cfg EncoderConfig) {
	enc.EncoderBase.SetConfig(cfg)

	if len(cfg.Fields) == 0 {
		enc.fieldsEncoded = ""

		return
	}

	buf := bytebufferpool.Get()
	buf.WriteString("{") // nolint:errcheck

	for i, field := range enc.cfg.Fields {
		if i > 0 {
			buf.WriteString(",") // nolint:errcheck
		}

		buf.WriteString("\"")      // nolint:errcheck
		buf.WriteString(field.Key) // nolint:errcheck
		buf.WriteString("\":\"")   // nolint:errcheck
		enc.WriteInterface(buf, field.Value)
		buf.WriteString("\"") // nolint:errcheck
	}

	buf.WriteString("}")     // nolint:errcheck
	buf.WriteString(sepText) // nolint:errcheck

	enc.SetFieldsEnconded(buf.String())

	bytebufferpool.Put(buf)
}

// Encode encodes the given level string, message and arguments to the buffer.
func (enc *EncoderText) Encode(buf *bytebufferpool.ByteBuffer, levelStr, msg string, args []interface{}) error {
	now := time.Now()
	if enc.cfg.UTC {
		now = now.UTC()
	}

	if enc.cfg.Datetime {
		enc.WriteDatetime(buf, now)
		buf.WriteString(sepText) // nolint:errcheck
	}

	if enc.cfg.Timestamp {
		enc.WriteTimestamp(buf, now)
		buf.WriteString(sepText) // nolint:errcheck
	}

	if levelStr != "" {
		buf.WriteString(levelStr) // nolint:errcheck
		buf.WriteString(sepText)  // nolint:errcheck
	}

	if enc.cfg.Shortfile || enc.cfg.Longfile {
		enc.WriteFileCaller(buf)
		buf.WriteString(sepText) // nolint:errcheck
	}

	enc.WriteFieldsEnconded(buf)

	enc.WriteMessage(buf, msg, args)

	enc.WriteNewLine(buf)

	return nil
}
