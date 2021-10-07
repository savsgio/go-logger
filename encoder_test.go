package logger

import (
	"fmt"
	"path"
	"reflect"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/valyala/bytebufferpool"
)

func newTestEncoderConfig() EncoderConfig {
	return EncoderConfig{}
}

func newTestEncoderBase() *EncoderBase {
	enc := new(EncoderBase)
	enc.SetConfig(newTestEncoderConfig())

	return enc
}

func TestEncoderBase_Config(t *testing.T) {
	enc := newTestEncoderBase()

	if cfg := enc.Config(); !reflect.DeepEqual(cfg, enc.cfg) {
		t.Errorf("cfg == %v, want %v", cfg, enc.cfg)
	}
}

func TestEncoderBase_SetConfig(t *testing.T) {
	cfg := EncoderConfig{
		Datetime:  true,
		calldepth: calldepth,
	}

	enc := newTestEncoderBase()
	enc.SetConfig(cfg)

	if !reflect.DeepEqual(enc.cfg, cfg) {
		t.Errorf("cfg == %v, want %v", enc.cfg, cfg)
	}
}

func TestEncoderBase_SetFieldsEnconded(t *testing.T) {
	fieldsEncoded := "v1 - v2 - v3"

	enc := newTestEncoderBase()
	enc.SetFieldsEnconded(fieldsEncoded)

	if enc.fieldsEncoded != fieldsEncoded {
		t.Errorf("fieldsEncoded == %s, want %s", enc.fieldsEncoded, fieldsEncoded)
	}
}

func TestEncoderBase_getFileCaller(t *testing.T) { // nolint:funlen
	type args struct {
		short     bool
		long      bool
		calldepth int
	}

	type want struct {
		file func(filepath string) string
		line func(line int) int
	}

	enc := newTestEncoderBase()

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "long",
			args: args{
				long:      true,
				short:     false,
				calldepth: 2,
			},
			want: want{
				file: func(filepath string) string {
					return filepath
				},
				line: func(line int) int {
					return line
				},
			},
		},
		{
			name: "short",
			args: args{
				long:      false,
				short:     true,
				calldepth: 2,
			},
			want: want{
				file: func(filepath string) string {
					_, wantFile := path.Split(filepath)

					return wantFile
				},
				line: func(line int) int {
					return line
				},
			},
		},
		{
			name: "invalid calldepth",
			args: args{
				long:      false,
				short:     true,
				calldepth: 1000,
			},
			want: want{
				file: func(filepath string) string {
					return "???"
				},
				line: func(line int) int {
					return 0
				},
			},
		},
	}

	for i := range tests {
		test := tests[i]

		t.Run(test.name, func(t *testing.T) {
			t.Helper()

			cfg := newTestEncoderConfig()
			cfg.Shortfile = test.args.short
			cfg.Longfile = test.args.long
			cfg.calldepth = test.args.calldepth

			enc.SetConfig(cfg)

			_, filepath, fileLine, _ := runtime.Caller(0)
			file, line := enc.getFileCaller()

			if wantFile := test.want.file(filepath); file != wantFile {
				t.Errorf("file == %s, want %s", file, wantFile)
			}

			if wantLine := test.want.line(fileLine + 1); line != wantLine {
				t.Errorf("line == %d, want %d", line, wantLine)
			}
		})
	}
}

func TestEncoderBase_WriteDatetime(t *testing.T) {
	enc := newTestEncoderBase()
	now := time.Now()

	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)

	enc.WriteDatetime(buf, now)

	wantDatetime := now.Format(time.RFC3339)

	if datetime := buf.String(); datetime != wantDatetime {
		t.Errorf("datetime == %s, want %s", datetime, wantDatetime)
	}
}

func TestEncoderBase_WriteTimestamp(t *testing.T) {
	enc := newTestEncoderBase()
	now := time.Now()

	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)

	enc.WriteTimestamp(buf, now)

	wantTs := strconv.FormatInt(now.Unix(), 10) // nolint:stylecheck

	if ts := buf.String(); ts != wantTs {
		t.Errorf("timestamp == %s, want %s", ts, wantTs)
	}
}

func TestEncoderBase_WriteFileCaller(t *testing.T) {
	cfg := newTestEncoderConfig()
	cfg.calldepth = 3

	enc := newTestEncoderBase()
	enc.SetConfig(cfg)

	getFileCaller := func() (string, int) {
		return enc.getFileCaller()
	}

	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)

	enc.WriteFileCaller(buf)

	file, line := getFileCaller()
	wantFileCaller := fmt.Sprintf("%s:%d", file, line-2)

	if fileCaller := buf.String(); fileCaller != wantFileCaller {
		t.Errorf("fileCaller == %s, want %s", fileCaller, wantFileCaller)
	}
}

func TestEncoderBase_WriteFieldsEnconded(t *testing.T) {
	enc := newTestEncoderBase()

	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)

	for _, wantFieldsEncoded := range []string{"", "v1 - v2"} {
		enc.SetFieldsEnconded(wantFieldsEncoded)

		enc.WriteFieldsEnconded(buf)

		if fieldsEncoded := buf.String(); fieldsEncoded != wantFieldsEncoded {
			t.Errorf("fieldsEncoded == %s, want %s", fieldsEncoded, wantFieldsEncoded)
		}

		buf.Reset()
	}
}

func TestEncoderBase_WriteMessage(t *testing.T) { // nolint:funlen
	type args struct {
		msg  string
		args []interface{}
	}

	type want struct {
		message string
	}

	enc := newTestEncoderBase()

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
				message: "Helloworld",
			},
		},
	}

	for i := range tests {
		test := tests[i]

		t.Run("", func(t *testing.T) {
			t.Helper()

			buf := bytebufferpool.Get()
			defer bytebufferpool.Put(buf)

			enc.WriteMessage(buf, test.args.msg, test.args.args)

			if message := buf.String(); message != test.want.message {
				t.Errorf("message == %s, want %s", message, test.want.message)
			}
		})
	}
}

func TestEncoderBase_WriteNewLine(t *testing.T) {
	enc := newTestEncoderBase()

	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)

	str := "foo"
	wantStr := str + "\n"

	buf.SetString(str)
	enc.WriteNewLine(buf)

	if bufStr := buf.String(); bufStr != wantStr {
		t.Errorf("line == %s, want %s", bufStr, wantStr)
	}

	buf.SetString(wantStr)
	enc.WriteNewLine(buf)

	if bufStr := buf.String(); bufStr != wantStr {
		t.Errorf("line == %s, want %s", bufStr, wantStr)
	}
}
