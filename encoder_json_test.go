package logger

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/valyala/bytebufferpool"
)

func newTestEncoderJSON() *EncoderJSON {
	enc := NewEncoderJSON()
	enc.SetConfig(newTestEncoderConfig())

	return enc
}

func Test_NewEncoderJSON(t *testing.T) {
	if enc := NewEncoderJSON(); enc == nil {
		t.Error("return nil")
	}
}

func TestEncoderJSON_hasBytesSpecialChars(t *testing.T) { // nolint:funlen
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
			enc := newTestEncoderJSON()

			if result := enc.hasBytesSpecialChars(test.args.value); result != test.want.result {
				t.Errorf("value == '%s', result = %t, want %t", test.args.value, result, test.want.result)
			}
		})
	}
}

func TestEncoderJSON_writeEscapedBytes(t *testing.T) { // nolint:funlen
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
			enc := newTestEncoderJSON()
			buf := bytebufferpool.Get()

			enc.writeEscapedBytes(buf, test.args.value)

			if result := buf.String(); result != test.want.result {
				t.Errorf("value == '%s', result = %s, want %s", test.args.value, result, test.want.result)
			}
		})
	}
}

func TestEncoderJSON_escape(t *testing.T) { // nolint:funlen
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
			enc := newTestEncoderJSON()

			buf := bytebufferpool.Get()
			buf.WriteString(line)      // nolint:errcheck
			buf.Write(test.args.value) // nolint:errcheck

			startAt := len(line)

			enc.escape(buf, startAt)

			if result := buf.String(); result != test.want.result {
				t.Errorf("value == '%s', result = %s, want %s", test.args.value, result, test.want.result)
			}
		})
	}
}

func TestEncoderJSON_WriteInterface(t *testing.T) {
	type args struct {
		value interface{}
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
				value: "hello world",
			},
			want: want{
				result: "hello world",
			},
		},
		{
			args: args{
				value: `hello \world"`,
			},
			want: want{
				result: `hello \\world\"`,
			},
		},
	}

	for i := range tests {
		test := tests[i]

		t.Run("", func(t *testing.T) {
			enc := newTestEncoderJSON()
			buf := bytebufferpool.Get()

			enc.WriteInterface(buf, test.args.value)

			if result := buf.String(); result != test.want.result {
				t.Errorf("value == '%v', result = %s, want %s", test.args.value, result, test.want.result)
			}
		})
	}
}

func TestEncoderJSON_WriteMessage(t *testing.T) {
	type args struct {
		msg  string
		args []interface{}
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
				msg:  `hello %s`,
				args: []interface{}{"world"},
			},
			want: want{
				result: "hello world",
			},
		},
		{
			args: args{
				msg:  `hello \%s"`,
				args: []interface{}{"world"},
			},
			want: want{
				result: `hello \\world\"`,
			},
		},
	}

	for i := range tests {
		test := tests[i]

		t.Run("", func(t *testing.T) {
			enc := newTestEncoderJSON()
			buf := bytebufferpool.Get()

			enc.WriteMessage(buf, test.args.msg, test.args.args)

			if result := buf.String(); result != test.want.result {
				t.Errorf(
					"msg == '%v' | args == %v, result = %s, want %s",
					test.args.msg, test.args.args, result, test.want.result,
				)
			}
		})
	}
}

func TestEncoderJSON_Copy(t *testing.T) {
	enc := newTestEncoderJSON()
	copyEnc, ok := enc.Copy().(*EncoderJSON)

	if !ok {
		t.Fatal("the copy is not a EncoderJSON pointer")
	}

	encPtr := reflect.ValueOf(enc).Pointer()
	copyEncPtr := reflect.ValueOf(copyEnc).Pointer()

	if copyEncPtr == encPtr {
		t.Error("the copy has the same pointer than original")
	}

	testEncoderBaseCopy(t, &enc.EncoderBase, &copyEnc.EncoderBase)
}

func TestEncoderJSON_SetConfig(t *testing.T) { // nolint:funlen
	type args struct {
		cfg EncoderConfig
	}

	type want struct {
		fieldsEncoded string
	}

	tests := []struct {
		args args
		want want
	}{
		{
			args: args{
				cfg: EncoderConfig{
					Timestamp: true,
					Flag:      Ltimestamp,
				},
			},
			want: want{
				fieldsEncoded: "",
			},
		},
		{
			args: args{
				cfg: EncoderConfig{
					Timestamp: true,
					Flag:      Ltimestamp,
					Fields:    []Field{{"foo", "bar"}, {"buzz", []int{1, 2, 3}}},
				},
			},
			want: want{
				fieldsEncoded: `"foo":"bar","buzz":"[1 2 3]",`,
			},
		},
		{
			args: args{
				cfg: EncoderConfig{
					Timestamp: true,
					Flag:      Ltimestamp,
					Fields:    []Field{{"foo", `id: "123"`}, {"buzz", []int{1, 2, 3}}},
				},
			},
			want: want{
				fieldsEncoded: `"foo":"id: \"123\"","buzz":"[1 2 3]",`,
			},
		},
	}

	for i := range tests {
		test := tests[i]

		t.Run("", func(t *testing.T) {
			t.Helper()

			enc := newTestEncoderJSON()
			enc.SetConfig(test.args.cfg)

			if encoderCfg := enc.Config(); !reflect.DeepEqual(encoderCfg, test.args.cfg) {
				t.Errorf("cfg == %v, want %v", encoderCfg, test.args.cfg)
			}

			if fieldsEncoded := enc.FieldsEnconded(); fieldsEncoded != test.want.fieldsEncoded {
				t.Errorf("fieldsEncoded == %s, want %s", fieldsEncoded, test.want.fieldsEncoded)
			}
		})
	}
}

func TestEncoderJSON_Encode(t *testing.T) { // nolint:funlen,dupl
	testCases := []testEncodeCase{
		{
			args: testEncodeArgs{
				cfg:      EncoderConfig{},
				levelStr: debugLevelStr,
				msg:      "Hello %s",
				args:     []interface{}{"world"},
			},
			want: testEncodeWant{
				lineRegexExpr: fmt.Sprintf(
					`^{"level":"%s","message":"%s"}\n$`,
					levelRegex, messageRegex,
				),
			},
		},
		{
			args: testEncodeArgs{
				cfg: EncoderConfig{
					UTC:       true,
					Datetime:  true,
					Timestamp: true,
					Shortfile: true,
					Longfile:  true,
				},
				levelStr: debugLevelStr,
				msg:      "Hello %s",
				args:     []interface{}{"world"},
			},
			want: testEncodeWant{
				lineRegexExpr: fmt.Sprintf(
					`^{"datetime":"%s","timestamp":"%s","level":"%s","file":"%s","message":"%s"}\n$`,
					datetimeRegex, timestampRegex, levelRegex, fileCallerRegex, messageRegex,
				),
			},
		},
		{
			args: testEncodeArgs{
				cfg: EncoderConfig{
					Fields:    []Field{{"foo", "bar"}},
					UTC:       true,
					Datetime:  true,
					Timestamp: true,
					Shortfile: true,
					Longfile:  true,
				},
				levelStr: debugLevelStr,
				msg:      "Hello %s",
				args:     []interface{}{"world"},
			},
			want: testEncodeWant{
				lineRegexExpr: fmt.Sprintf(
					`^{"datetime":"%s","timestamp":"%s","level":"%s","file":"%s",%s,"message":"%s"}\n$`,
					datetimeRegex, timestampRegex, levelRegex, fileCallerRegex, fieldsKVRegex, messageRegex,
				),
			},
		},
		{
			args: testEncodeArgs{
				cfg: EncoderConfig{
					Fields:    []Field{{"foo", `id: "bar"`}},
					UTC:       true,
					Datetime:  true,
					Timestamp: true,
					Shortfile: true,
					Longfile:  true,
				},
				levelStr: debugLevelStr,
				msg:      "Hello %s",
				args:     []interface{}{"world"},
			},
			want: testEncodeWant{
				lineRegexExpr: fmt.Sprintf(
					`^{"datetime":"%s","timestamp":"%s","level":"%s","file":"%s",%s,"message":"%s"}\n$`,
					datetimeRegex, timestampRegex, levelRegex, fileCallerRegex, fieldsKVRegex, messageRegex,
				),
			},
		},
		{
			args: testEncodeArgs{
				cfg: EncoderConfig{
					Fields:    []Field{{"foo", "bar"}},
					UTC:       true,
					Datetime:  true,
					Timestamp: true,
					Shortfile: true,
					Longfile:  true,
				},
				levelStr: printLevelStr,
				msg:      "Hello %s",
				args:     []interface{}{"world"},
			},
			want: testEncodeWant{
				lineRegexExpr: fmt.Sprintf(
					`^{"datetime":"%s","timestamp":"%s","file":"%s",%s,"message":"%s"}\n$`,
					datetimeRegex, timestampRegex, fileCallerRegex, fieldsKVRegex, messageRegex,
				),
			},
		},
	}

	enc := newTestEncoderJSON()

	testEncoderEncode(t, enc, testCases)
}

func BenchmarkEncoderJSON_Encode(b *testing.B) {
	enc := newTestEncoderJSON()
	benchmarkEncoderEncode(b, enc)
}
