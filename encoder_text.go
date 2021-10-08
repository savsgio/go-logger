package logger

import (
	"time"

	"github.com/valyala/bytebufferpool"
)

const sepText = " - "

func NewEncoderText() *EncoderText {
	return new(EncoderText)
}

func (enc *EncoderText) Copy() Encoder {
	copyEnc := NewEncoderText()
	copyEnc.EncoderBase = *enc.EncoderBase.Copy()

	return copyEnc
}

func (enc *EncoderText) SetConfig(cfg EncoderConfig) {
	enc.EncoderBase.SetConfig(cfg)

	if len(cfg.Fields) == 0 {
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
