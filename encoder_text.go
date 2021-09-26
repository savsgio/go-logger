package logger

import (
	"time"

	"github.com/valyala/bytebufferpool"
)

func NewEncoderText() *EncoderText {
	return new(EncoderText)
}

func (enc *EncoderText) Encode(level, msg string, args []interface{}) error {
	buf := bytebufferpool.Get()

	now := time.Now()
	if enc.opts.UTC {
		now = now.UTC()
	}

	if enc.opts.Date || enc.opts.Time {
		if enc.opts.Date {
			enc.WriteDate(buf, now)
		}

		if enc.opts.Time {
			if enc.opts.Date {
				buf.WriteString(" ") // nolint:errcheck
			}

			enc.WriteTime(buf, now, enc.opts.TimeMicroseconds)
		}

		buf.WriteString(" - ") // nolint:errcheck
	}

	if enc.opts.Shortfile || enc.opts.Longfile {
		enc.WriteFileCaller(buf)
		buf.WriteString(" - ") // nolint:errcheck
	}

	buf.WriteString(level) // nolint:errcheck
	buf.WriteString(" - ") // nolint:errcheck
	enc.WriteMessage(buf, msg, args)
	enc.WriteNewLine(buf)

	_, err := enc.Write(buf.Bytes())

	bytebufferpool.Put(buf)

	return err
}
