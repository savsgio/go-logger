package logger

import (
	"time"

	"github.com/valyala/bytebufferpool"
)

func NewEncoderJSON() *EncoderJSON {
	return new(EncoderJSON)
}

func (enc *EncoderJSON) Encode(level, msg string, args []interface{}) error {
	now := time.Now()
	if enc.opts.UTC {
		now = now.UTC()
	}

	buf := bytebufferpool.Get()
	buf.WriteByte('{') // nolint:errcheck

	if enc.opts.Date || enc.opts.Time {
		if enc.opts.Date {
			buf.WriteString("\"time\":\"") // nolint:errcheck
			enc.WriteDate(buf, now)
		}

		if enc.opts.Time {
			if enc.opts.Date {
				buf.WriteString(" ") // nolint:errcheck
			}

			enc.WriteTime(buf, now, enc.opts.TimeMicroseconds)
		}

		buf.WriteString("\",") // nolint:errcheck
	}

	if enc.opts.Shortfile || enc.opts.Longfile {
		buf.WriteString("\"file\":\"") // nolint:errcheck
		enc.WriteFileCaller(buf)
		buf.WriteString("\",") // nolint:errcheck
	}

	buf.WriteString("\"level\":\"")  // nolint:errcheck
	buf.WriteString(level)           // nolint:errcheck
	buf.WriteString("\",\"msg\":\"") // nolint:errcheck
	enc.WriteMessage(buf, msg, args)
	buf.WriteString("\"}") // nolint:errcheck

	enc.WriteNewLine(buf)

	_, err := enc.Write(buf.Bytes())

	bytebufferpool.Put(buf)

	return err
}