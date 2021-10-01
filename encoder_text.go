package logger

import (
	"time"

	"github.com/valyala/bytebufferpool"
)

func NewEncoderText() *EncoderText {
	return new(EncoderText)
}

func (enc *EncoderText) Copy() Encoder {
	copyEnc := NewEncoderText()
	copyEnc.cfg = enc.cfg

	return copyEnc
}

func (enc *EncoderText) SetConfig(cfg Config) {
	enc.EncoderBase.SetConfig(cfg)

	buf := bytebufferpool.Get()

	for _, field := range enc.cfg.Fields {
		buf.WriteString(field.Key) // nolint:errcheck
		buf.WriteString(": ")      // nolint:errcheck
		enc.WriteInterface(buf, field.Value)
		buf.WriteString(" - ") // nolint:errcheck
	}

	enc.SetFieldsEnconded(buf.String())

	bytebufferpool.Put(buf)
}

func (enc *EncoderText) Encode(level, msg string, args []interface{}) error {
	buf := bytebufferpool.Get()

	now := time.Now()
	if enc.cfg.UTC {
		now = now.UTC()
	}

	if enc.cfg.Datetime {
		buf.WriteString("datetime: ") // nolint:errcheck
		enc.WriteDatetime(buf, now)
		buf.WriteString(" - ") // nolint:errcheck
	}

	if enc.cfg.Timestamp {
		buf.WriteString("timestamp: ") // nolint:errcheck
		enc.WriteTimestamp(buf, now)
		buf.WriteString(" - ") // nolint:errcheck
	}

	if enc.cfg.Shortfile || enc.cfg.Longfile {
		buf.WriteString("file: ") // nolint:errcheck
		enc.WriteFileCaller(buf)
		buf.WriteString(" - ") // nolint:errcheck
	}

	if level != "" {
		buf.WriteString("level: ") // nolint:errcheck
		buf.WriteString(level)     // nolint:errcheck
		buf.WriteString(" - ")     // nolint:errcheck
	}

	enc.WriteFieldsEnconded(buf)

	buf.WriteString("message: ") // nolint:errcheck
	enc.WriteMessage(buf, msg, args)

	enc.WriteNewLine(buf)

	_, err := enc.Write(buf.Bytes())

	bytebufferpool.Put(buf)

	return err
}
