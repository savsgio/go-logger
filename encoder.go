package logger

import (
	"fmt"
	"io"
	"runtime"
	"strconv"
	"time"

	"github.com/valyala/bytebufferpool"
)

func (enc *EncoderBase) SetOutput(output io.Writer) {
	enc.output = output
}

func (enc *EncoderBase) SetOptions(opts Options) {
	enc.opts = opts
}

func (enc *EncoderBase) GetCaller() (string, int) {
	pc := make([]uintptr, 1)

	numFrames := runtime.Callers(calldepth+1, pc)
	if numFrames < 1 {
		return "???", 0
	}

	frame, _ := runtime.CallersFrames(pc).Next()
	file := frame.File

	if enc.opts.Shortfile {
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				file = file[i+1:]

				break
			}
		}
	}

	return file, frame.Line
}

func (enc *EncoderBase) WritePadInt(buf *bytebufferpool.ByteBuffer, i int, wid int) { // nolint:interfacer
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1

	for i >= 10 || wid > 1 {
		wid--

		q := i / 10
		b[bp] = byte('0' + i - q*10)

		bp--

		i = q
	}

	// i < 10
	b[bp] = byte('0' + i)

	buf.Write(b[bp:]) // nolint:errcheck
}

func (enc *EncoderBase) WriteDate(buf *bytebufferpool.ByteBuffer, now time.Time) {
	year, month, day := now.Date()

	enc.WritePadInt(buf, year, 4)
	buf.WriteByte('/') // nolint:errcheck

	enc.WritePadInt(buf, int(month), 2)
	buf.WriteByte('/') // nolint:errcheck

	enc.WritePadInt(buf, day, 2)
}

func (enc *EncoderBase) WriteTime(buf *bytebufferpool.ByteBuffer, now time.Time, withMicroseconds bool) {
	hour, min, sec := now.Clock()

	enc.WritePadInt(buf, hour, 2)
	buf.WriteByte(':') // nolint:errcheck

	enc.WritePadInt(buf, min, 2)
	buf.WriteByte(':') // nolint:errcheck

	enc.WritePadInt(buf, sec, 2)

	if withMicroseconds {
		buf.WriteByte('.') // nolint:errcheck
		enc.WritePadInt(buf, now.Nanosecond()/1e3, 6)
	}
}

func (enc *EncoderBase) WriteFileCaller(buf *bytebufferpool.ByteBuffer) {
	file, line := enc.GetCaller()

	buf.WriteString(file) // nolint:errcheck
	buf.WriteByte(':')    // nolint:errcheck
	buf.B = strconv.AppendInt(buf.B, int64(line), 10)
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
		}
	default:
		fmt.Fprint(buf, args...)
	}
}

func (enc *EncoderBase) WriteNewLine(buf *bytebufferpool.ByteBuffer) {
	if buf.B[buf.Len()-1] != '\n' {
		buf.WriteByte('\n') // nolint:errcheck
	}
}

func (enc *EncoderBase) Write(p []byte) (int, error) {
	return enc.output.Write(p) // nolint:wrapcheck
}
