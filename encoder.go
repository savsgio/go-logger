package logger

import (
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/valyala/bytebufferpool"
)

func (enc *EncoderBase) Config() EncoderConfig {
	return enc.cfg
}

func (enc *EncoderBase) SetConfig(cfg EncoderConfig) {
	enc.cfg = cfg
}

func (enc *EncoderBase) SetFieldsEnconded(fieldsEncoded string) {
	enc.fieldsEncoded = fieldsEncoded
}

func (enc *EncoderBase) getFileCaller() (string, int) {
	pc := make([]uintptr, 1)

	numFrames := runtime.Callers(enc.cfg.calldepth, pc)
	if numFrames < 1 {
		return "???", 0
	}

	frame, _ := runtime.CallersFrames(pc).Next()
	file := frame.File

	if enc.cfg.Shortfile {
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				file = file[i+1:]

				break
			}
		}
	}

	return file, frame.Line
}

func (enc *EncoderBase) WriteDatetime(buf *bytebufferpool.ByteBuffer, now time.Time) {
	buf.B = now.AppendFormat(buf.B, time.RFC3339)
}

func (enc *EncoderBase) WriteTimestamp(buf *bytebufferpool.ByteBuffer, now time.Time) {
	buf.B = strconv.AppendInt(buf.B, now.Unix(), 10)
}

func (enc *EncoderBase) WriteFileCaller(buf *bytebufferpool.ByteBuffer) {
	file, line := enc.getFileCaller()

	buf.WriteString(file) // nolint:errcheck
	buf.WriteByte(':')    // nolint:errcheck
	buf.B = strconv.AppendInt(buf.B, int64(line), 10)
}

func (enc *EncoderBase) WriteFieldsEnconded(buf *bytebufferpool.ByteBuffer) { // nolint:interfacer
	if enc.fieldsEncoded != "" {
		buf.WriteString(enc.fieldsEncoded) // nolint:errcheck
	}
}

func (enc *EncoderBase) WriteInterface(buf *bytebufferpool.ByteBuffer, value interface{}) {
	if str, ok := value.(string); ok {
		buf.WriteString(str) // nolint:errcheck
	} else {
		fmt.Fprint(buf, value)
	}
}

func (enc *EncoderBase) WriteMessage(buf *bytebufferpool.ByteBuffer, msg string, args []interface{}) {
	lenArgs := len(args)

	switch {
	case lenArgs == 0:
		buf.WriteString(msg) // nolint:errcheck
	case msg != "":
		fmt.Fprintf(buf, msg, args...)
	case lenArgs == 1:
		if str, ok := args[0].(string); ok {
			buf.WriteString(str) // nolint:errcheck

			return
		}

		fallthrough
	default:
		fmt.Fprint(buf, args...)
	}
}

func (enc *EncoderBase) WriteNewLine(buf *bytebufferpool.ByteBuffer) {
	if buf.B[buf.Len()-1] != '\n' {
		buf.WriteByte('\n') // nolint:errcheck
	}
}
