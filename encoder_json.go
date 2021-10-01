package logger

import (
	"time"

	"github.com/valyala/bytebufferpool"
)

func NewEncoderJSON() *EncoderJSON {
	return new(EncoderJSON)
}

func (enc *EncoderJSON) Copy() Encoder {
	copyEnc := NewEncoderJSON()
	copyEnc.cfg = enc.cfg

	return copyEnc
}

func (enc *EncoderJSON) SetConfig(cfg Config) {
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

func (enc *EncoderJSON) Encode(level, msg string, args []interface{}) error {
	now := time.Now()
	if enc.cfg.UTC {
		now = now.UTC()
	}

	buf := bytebufferpool.Get()
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

	if enc.cfg.Shortfile || enc.cfg.Longfile {
		buf.WriteString("\"file\":\"") // nolint:errcheck
		enc.WriteFileCaller(buf)
		buf.WriteString("\",") // nolint:errcheck
	}

	if level != "" {
		buf.WriteString("\"level\":\"") // nolint:errcheck
		buf.WriteString(level)          // nolint:errcheck
		buf.WriteString("\",")          // nolint:errcheck
	}

	enc.WriteFieldsEnconded(buf)

	buf.WriteString("\"message\":\"") // nolint:errcheck
	enc.WriteMessage(buf, msg, args)
	buf.WriteString("\"}") // nolint:errcheck

	enc.WriteNewLine(buf)

	_, err := enc.Write(buf.Bytes())

	bytebufferpool.Put(buf)

	return err
}
