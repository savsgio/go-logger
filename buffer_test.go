package logger

import (
	"bytes"
	"fmt"
	"path"
	"strconv"
	"testing"
	"time"
)

func TestBuffer_hasBytesSpecialChars(t *testing.T) { // nolint:funlen
	type args struct {
		value []byte
	}

	type want struct {
		result bool
	}

	tests := []struct {
		args args
		want want
	}{
		{
			args: args{
				value: []byte("some string"),
			},
			want: want{
				result: false,
			},
		},
		{
			args: args{
				value: []byte(`"some" string`),
			},
			want: want{
				result: true,
			},
		},
		{
			args: args{
				value: []byte(`some \string\`),
			},
			want: want{
				result: true,
			},
		},
		{
			args: args{
				value: []byte{0x23, 0x56, 0x10, 0x67},
			},
			want: want{
				result: true,
			},
		},
		{
			args: args{
				value: []byte{0x23, 0x56, 0x20, 0x67},
			},
			want: want{
				result: false,
			},
		},
	}

	for i := range tests {
		test := tests[i]

		t.Run("", func(t *testing.T) {
			buf := NewBuffer()

			if result := buf.hasBytesSpecialChars(test.args.value); result != test.want.result {
				t.Errorf("value == '%s', result = %t, want %t", test.args.value, result, test.want.result)
			}
		})
	}
}

func TestBuffer_writeEscapedBytes(t *testing.T) { // nolint:funlen
	type args struct {
		value []byte
	}

	type want struct {
		result string
	}

	tests := []struct {
		args args
		want want
	}{
		{
			args: args{
				value: []byte("some string"),
			},
			want: want{
				result: "some string",
			},
		},
		{
			args: args{
				value: []byte(`"some" string`),
			},
			want: want{
				result: `\"some\" string`,
			},
		},
		{
			args: args{
				value: []byte(`some: \string\`),
			},
			want: want{
				result: `some: \\string\\`,
			},
		},
	}

	for i := range tests {
		test := tests[i]

		t.Run("", func(t *testing.T) {
			buf := NewBuffer()
			buf.writeEscapedBytes(test.args.value)

			if result := buf.String(); result != test.want.result {
				t.Errorf("value == '%s', result = %s, want %s", test.args.value, result, test.want.result)
			}
		})
	}
}

func TestBuffer_formatMessage(t *testing.T) { // nolint:funlen
	type args struct {
		msg  string
		args []interface{}
	}

	type want struct {
		message string
	}

	tests := []struct {
		args args
		want want
	}{
		{
			args: args{
				msg: "Hello world",
			},
			want: want{
				message: "Hello world",
			},
		},
		{
			args: args{
				msg:  "Hello %s",
				args: []interface{}{"world"},
			},
			want: want{
				message: "Hello world",
			},
		},
		{
			args: args{
				msg:  "",
				args: []interface{}{"Hello world"},
			},
			want: want{
				message: "Hello world",
			},
		},
		{
			args: args{
				msg:  "",
				args: []interface{}{1}, // case: fallthrough
			},
			want: want{
				message: "1",
			},
		},
		{
			args: args{
				args: []interface{}{"Hello", "world"},
			},
			want: want{
				message: fmt.Sprint("Hello", "world"),
			},
		},
	}

	for i := range tests {
		test := tests[i]

		t.Run("", func(t *testing.T) {
			buf := NewBuffer()

			message := buf.formatMessage(test.args.msg, test.args.args)

			if message != test.want.message {
				t.Errorf("message == %s, want %s", message, test.want.message)
			}
		})
	}
}

func TestBuffer_Reset(t *testing.T) {
	buf := NewBuffer()

	buf.WriteString("hello world") // nolint:errcheck
	buf.formatMessage("yeah %s", []interface{}{"sir"})

	buf.Reset()

	if result := buf.b1.Len(); result > 0 {
		t.Errorf("Buffer.b1 length is not zero")
	}

	if result := buf.b2.Len(); result > 0 {
		t.Errorf("Buffer.b2 length is not zero")
	}
}

func TestBuffer_Len(t *testing.T) {
	buf := NewBuffer()

	buf.WriteString("hello world") // nolint:errcheck
	buf.formatMessage("yeah %s", []interface{}{"sir"})

	wantLen := buf.b1.Len()

	if result := buf.Len(); result != wantLen {
		t.Errorf("result == %d, want %d", result, wantLen)
	}
}

func TestBuffer_String(t *testing.T) {
	buf := NewBuffer()
	wantString := "hello sir"

	buf.WriteString(wantString) // nolint:errcheck

	if result := buf.String(); result != wantString {
		t.Errorf("result == %s, want %s", result, wantString)
	}
}

func TestBuffer_Bytes(t *testing.T) {
	buf := NewBuffer()
	wantString := "hello men"

	buf.WriteString(wantString) // nolint:errcheck

	if result := buf.Bytes(); string(result) != wantString {
		t.Errorf("result == %s, want %s", result, wantString)
	}
}

func TestBuffer_Escape(t *testing.T) { // nolint:funlen
	type args struct {
		value []byte
	}

	type want struct {
		result string
	}

	line := "test line"

	tests := []struct {
		args args
		want want
	}{
		{
			args: args{
				value: []byte("some string"),
			},
			want: want{
				result: line + "some string",
			},
		},
		{
			args: args{
				value: []byte(`"some" string`),
			},
			want: want{
				result: line + `\"some\" string`,
			},
		},
		{
			args: args{
				value: []byte(`some: \string\`),
			},
			want: want{
				result: line + `some: \\string\\`,
			},
		},
	}

	for i := range tests {
		test := tests[i]

		t.Run("", func(t *testing.T) {
			startAt := len(line)

			buf := NewBuffer()
			buf.WriteString(line)      // nolint:errcheck
			buf.Write(test.args.value) // nolint:errcheck

			buf.Escape(startAt)

			if result := buf.String(); result != test.want.result {
				t.Errorf("value == '%s', result = %s, want %s", test.args.value, result, test.want.result)
			}
		})
	}
}

func TestBuffer_Write(t *testing.T) {
	buf := NewBuffer()
	wantValue := []byte("hello world")

	buf.Write(wantValue) // nolint:errcheck

	if result := buf.b1.Bytes(); string(result) != string(wantValue) {
		t.Errorf("result == %s, want %s", result, wantValue)
	}
}

func TestBuffer_WriteByte(t *testing.T) {
	buf := NewBuffer()
	wantValue := byte('b')

	buf.WriteByte(wantValue) // nolint:errcheck

	if result := buf.b1.Bytes(); string(result) != string(wantValue) {
		t.Errorf("result == %v, want %v", result, wantValue)
	}
}

func TestBuffer_WriteString(t *testing.T) {
	buf := NewBuffer()
	wantValue := "hello world"

	buf.WriteString(wantValue) // nolint:errcheck

	if result := buf.b1.String(); result != wantValue {
		t.Errorf("result == %s, want %s", result, wantValue)
	}
}

func TestBuffer_WriteTo(t *testing.T) {
	buf := NewBuffer()
	wantValue := "hello world"

	dst := new(bytes.Buffer)

	buf.WriteString(wantValue) // nolint:errcheck
	buf.WriteTo(dst)           // nolint:errcheck

	if result := dst.String(); result != wantValue {
		t.Errorf("result == %s, want %s", result, wantValue)
	}
}

func TestBuffer_WriteDatetime(t *testing.T) {
	buf := NewBuffer()
	now := time.Now()

	buf.WriteDatetime(now)

	wantDatetime := now.Format(time.RFC3339)

	if datetime := buf.String(); datetime != wantDatetime {
		t.Errorf("datetime == %s, want %s", datetime, wantDatetime)
	}
}

func TestBuffer_WriteTimestamp(t *testing.T) {
	buf := NewBuffer()
	now := time.Now()

	buf.WriteTimestamp(now)

	wantTs := strconv.FormatInt(now.Unix(), 10) // nolint:stylecheck

	if ts := buf.String(); ts != wantTs {
		t.Errorf("timestamp == %s, want %s", ts, wantTs)
	}
}

func TestBuffer_WriteFileCaller(t *testing.T) {
	caller := getFileCaller(2)

	// Short
	t.Run("Short", func(t *testing.T) {
		buf := NewBuffer()
		buf.WriteFileCaller(caller, true)

		_, filename := path.Split(caller.File)

		wantFileCaller := fmt.Sprintf("%s:%d", filename, caller.Line)

		if fileCaller := buf.String(); fileCaller != wantFileCaller {
			t.Errorf("fileCaller (short) == %s, want %s", fileCaller, wantFileCaller)
		}
	})

	t.Run("Long", func(t *testing.T) {
		buf := NewBuffer()
		buf.WriteFileCaller(caller, false)

		wantFileCaller := fmt.Sprintf("%s:%d", caller.File, caller.Line)

		if fileCaller := buf.String(); fileCaller != wantFileCaller {
			t.Errorf("fileCaller (long) == %s, want %s", fileCaller, wantFileCaller)
		}
	})
}

func TestBuffer_WriteInterface(t *testing.T) {
	buf := NewBuffer()
	value := []int{1, 2, 3}

	buf.WriteInterface(value)

	wantValue := fmt.Sprint(value)

	if result := buf.String(); result != wantValue {
		t.Errorf("value == %s, want %s", result, wantValue)
	}
}

func TestBuffer_WriteNewLine(t *testing.T) {
	buf := NewBuffer()

	str := "foo"
	wantStr := str + "\n"

	buf.WriteString(str) // nolint:errcheck
	buf.WriteNewLine()

	if bufStr := buf.String(); bufStr != wantStr {
		t.Errorf("line == %s, want %s", bufStr, wantStr)
	}

	buf.Reset()
	buf.WriteString(wantStr) // nolint:errcheck
	buf.WriteNewLine()

	if bufStr := buf.String(); bufStr != wantStr {
		t.Errorf("line == %s, want %s", bufStr, wantStr)
	}
}
