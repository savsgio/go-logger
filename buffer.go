package logger

import (
	"bytes"
	"fmt"
	"io"
	"runtime"
	"strconv"
	"sync"
	"time"

	gstrconv "github.com/savsgio/gotils/strconv"
	"github.com/valyala/bytebufferpool"
)

var bufferPool = sync.Pool{
	New: func() interface{} {
		return NewBuffer()
	},
}

// NewBuffer returns a new buffer.
func NewBuffer() *Buffer {
	return new(Buffer)
}

// AcquireBuffer returns a buffer instance from pool.
func AcquireBuffer() *Buffer {
	return bufferPool.Get().(*Buffer) // nolint:forcetypeassert
}

// ReleaseBuffer puts the given buffer in the pool after it is reset.
func ReleaseBuffer(b *Buffer) {
	b.Reset()
	bufferPool.Put(b)
}

func (b *Buffer) hasBytesSpecialChars(value []byte) bool {
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

func (b *Buffer) writeEscapedBytes(value []byte) {
	str := bytebufferpool.Get()
	str.Set(value) // NOTE: Use as a copy of b.

	str.B = strconv.AppendQuote(str.B, gstrconv.B2S(str.B))

	b.Write(str.B[len(value)+1 : str.Len()-1]) // nolint:errcheck
	bytebufferpool.Put(str)
}

func (b *Buffer) formatMessage(msg string, args []interface{}) string {
	b.b2.Reset()

	lenArgs := len(args)

	switch {
	case lenArgs == 0:
		b.b2.WriteString(msg) // nolint:errcheck
	case msg != "":
		fmt.Fprintf(&b.b2, msg, args...)
	case lenArgs == 1:
		if strValue, ok := args[0].(string); ok {
			b.b2.WriteString(strValue) // nolint:errcheck

			return b.b2.String()
		}

		fallthrough
	default:
		fmt.Fprint(&b.b2, args...)
	}

	return b.b2.String()
}

// Reset clears the buffer.
func (b *Buffer) Reset() {
	b.b1.Reset()
	b.b2.Reset()
}

// Len returns the size of the buffer.
func (b *Buffer) Len() int {
	return b.b1.Len()
}

// String returns string representation.
func (b *Buffer) String() string {
	return b.b1.String()
}

// Bytes returns all the accumulated bytes.
func (b *Buffer) Bytes() []byte {
	return b.b1.Bytes()
}

// Escape escapes accumulated bytes since the given index.
func (b *Buffer) Escape(startAt int) {
	if value := b.b1.B[startAt:]; b.hasBytesSpecialChars(value) {
		b.b1.Set(b.b1.B[:startAt])
		b.writeEscapedBytes(value)
	}
}

// Write writes the given bytes slice to the buffer.
func (b *Buffer) Write(s []byte) (int, error) {
	return b.b1.Write(s) // nolint:wrapcheck
}

// WriteByte writes the given byte to the buffer.
func (b *Buffer) WriteByte(s byte) error {
	return b.b1.WriteByte(s) // nolint:wrapcheck
}

// WriteString writes the given string to the buffer.
func (b *Buffer) WriteString(s string) (int, error) {
	return b.b1.WriteString(s) // nolint:wrapcheck
}

// WriteTo implements io.WriterTo.
func (b *Buffer) WriteTo(w io.Writer) (int64, error) {
	return b.b1.WriteTo(w) // nolint:wrapcheck
}

// WriteDatetime writes the given time to the buffer formatted with the given layout.
func (b *Buffer) WriteDatetime(now time.Time, layout string) {
	b.b1.B = now.AppendFormat(b.b1.B, layout)
}

// WriteTimestamp writes the timestamp to the buffer from the given time.
func (b *Buffer) WriteTimestamp(now time.Time, format TimestampFormat) {
	fn := now.Unix
	if format == TimestampFormatNanoseconds {
		fn = now.UnixNano
	}

	b.b1.B = strconv.AppendInt(b.b1.B, fn(), 10)
}

// WriteFileCaller writes the file caller to the buffer.
func (b *Buffer) WriteFileCaller(f runtime.Frame, short bool) {
	file := f.File

	if short {
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				file = file[i+1:]

				break
			}
		}
	}

	b.WriteString(file) // nolint:errcheck
	b.WriteByte(':')    // nolint:errcheck
	b.b1.B = strconv.AppendInt(b.b1.B, int64(f.Line), 10)
}

// WriteInterface writes an interface value to the buffer.
func (b *Buffer) WriteInterface(value interface{}) {
	if strValue, ok := value.(string); ok {
		b.WriteString(strValue) // nolint:errcheck
	} else {
		fmt.Fprint(b, value)
	}
}

// WriteNewLine writes a new line to the buffer if it's needed.
func (b *Buffer) WriteNewLine() {
	if length := b.Len(); length > 0 && b.b1.B[length-1] != '\n' {
		b.WriteByte('\n') // nolint:errcheck
	}
}
