package logger

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func newTestEncoderJSON() *EncoderJSON {
	cfg := newTestConfig()

	enc := NewEncoderJSON(EncoderJSONConfig{})
	enc.Configure(cfg)

	return enc
}

func Test_NewEncoderJSON(t *testing.T) { // nolint:funlen
	type args struct {
		cfg EncoderJSONConfig
	}

	type want struct {
		cfg EncoderJSONConfig
	}

	tests := []struct {
		args args
		want want
	}{
		{
			args: args{
				cfg: EncoderJSONConfig{},
			},
			want: want{
				cfg: EncoderJSONConfig{
					FieldMap: EnconderJSONFieldMap{
						DatetimeKey:  defaultJSONFieldKeyDatetime,
						TimestampKey: defaultJSONFieldKeyTimestamp,
						LevelKey:     defaultJSONFieldKeyLevel,
						FileKey:      defaultJSONFieldKeyFile,
						FunctionKey:  defaultJSONFieldKeyFunction,
						MessageKey:   defaultJSONFieldKeyMessage,
					},
					DatetimeLayout:  defaultDatetimeLayout,
					TimestampFormat: defaultTimestampFormat,
				},
			},
		},
		{
			args: args{
				cfg: EncoderJSONConfig{
					FieldMap: EnconderJSONFieldMap{
						DatetimeKey:  "@date",
						TimestampKey: "@time",
						LevelKey:     "log.level",
						FileKey:      "caller.file",
						FunctionKey:  "caller.func",
						MessageKey:   "msg",
					},
					DatetimeLayout:  time.RFC1123,
					TimestampFormat: TimestampFormatNanoseconds,
				},
			},
			want: want{
				cfg: EncoderJSONConfig{
					FieldMap: EnconderJSONFieldMap{
						DatetimeKey:  "@date",
						TimestampKey: "@time",
						LevelKey:     "log.level",
						FileKey:      "caller.file",
						FunctionKey:  "caller.func",
						MessageKey:   "msg",
					},
					DatetimeLayout:  time.RFC1123,
					TimestampFormat: TimestampFormatNanoseconds,
				},
			},
		},
	}

	for i := range tests {
		test := tests[i]

		t.Run("", func(t *testing.T) {
			enc := NewEncoderJSON(test.args.cfg)
			if enc == nil {
				t.Fatal("return nil")
			}

			if !reflect.DeepEqual(enc.cfg, test.want.cfg) {
				t.Errorf("confg == %v, want %v", enc.cfg, test.want.cfg)
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

func TestEncoderJSON_keys(t *testing.T) {
	type args struct {
		cfg Config
	}

	type want struct {
		result []string
	}

	tests := []struct {
		args args
		want want
	}{
		{
			args: args{
				cfg: Config{},
			},
			want: want{
				result: []string{
					defaultJSONFieldKeyLevel,
					defaultJSONFieldKeyMessage,
				},
			},
		},
		{
			args: args{
				cfg: Config{
					Datetime:  true,
					Timestamp: true,
					UTC:       true,
					Shortfile: true,
					Longfile:  true,
					Function:  true,
				},
			},
			want: want{
				result: []string{
					defaultJSONFieldKeyDatetime,
					defaultJSONFieldKeyTimestamp,
					defaultJSONFieldKeyLevel,
					defaultJSONFieldKeyFile,
					defaultJSONFieldKeyFunction,
					defaultJSONFieldKeyMessage,
				},
			},
		},
	}

	enc := newTestEncoderJSON()

	for i := range tests {
		test := tests[i]

		t.Run("", func(t *testing.T) {
			result := enc.keys(test.args.cfg)

			if !reflect.DeepEqual(result, test.want.result) {
				t.Errorf("keys == %v, want %v", result, test.want.result)
			}
		})
	}
}

func TestEncoderJSON_Configure(t *testing.T) { // nolint:funlen
	type args struct {
		cfg Config
	}

	type want struct {
		fieldsEncoded string
	}

	enc := newTestEncoderJSON()

	tests := []struct {
		args args
		want want
	}{
		{
			args: args{
				cfg: Config{},
			},
			want: want{
				fieldsEncoded: "",
			},
		},
		{
			args: args{
				cfg: Config{
					Fields: []Field{
						{enc.cfg.FieldMap.DatetimeKey, "hello"}, {"foo", "bar"}, {"buzz", []int{1, 2, 3}},
					},
					Datetime: false,
				},
			},
			want: want{
				fieldsEncoded: `"` + enc.cfg.FieldMap.DatetimeKey + `":"hello","foo":"bar","buzz":"[1 2 3]",`,
			},
		},
		{
			args: args{
				cfg: Config{
					Fields: []Field{
						{enc.cfg.FieldMap.DatetimeKey, "hello"}, {"foo", "bar"}, {"buzz", []int{1, 2, 3}},
					},
					Datetime: true,
				},
			},
			want: want{
				fieldsEncoded: `"fields.` + enc.cfg.FieldMap.DatetimeKey + `":"hello","foo":"bar","buzz":"[1 2 3]",`,
			},
		},
		{
			args: args{
				cfg: Config{
					Fields: []Field{{"foo", `id: "123"`}, {"buzz", []int{1, 2, 3}}},
				},
			},
			want: want{
				fieldsEncoded: `"foo":"id: \"123\"","buzz":"[1 2 3]",`,
			},
		},
		{
			args: args{
				cfg: Config{
					Fields: []Field{{`foo"ter"`, `id: "123"`}, {"buzz", []int{1, 2, 3}}},
				},
			},
			want: want{
				fieldsEncoded: `"foo\"ter\"":"id: \"123\"","buzz":"[1 2 3]",`,
			},
		},
	}

	for i := range tests {
		test := tests[i]

		t.Run("", func(t *testing.T) {
			t.Helper()

			enc.Configure(test.args.cfg)

			if fieldsEncoded := enc.FieldsEncoded(); fieldsEncoded != test.want.fieldsEncoded {
				t.Errorf("fieldsEncoded == %s, want %s", fieldsEncoded, test.want.fieldsEncoded)
			}
		})
	}
}

func TestEncoderJSON_Encode(t *testing.T) { // nolint:funlen,dupl
	testCases := []testEncodeCase{
		{
			args: testEncodeArgs{
				cfg:   Config{},
				level: DEBUG,
				msg:   "Hello %s",
				args:  []interface{}{"world"},
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
				cfg: Config{
					UTC:       true,
					Datetime:  true,
					Timestamp: true,
					Shortfile: true,
					Longfile:  true,
					Function:  true,
				},
				level: DEBUG,
				msg:   "Hello %s",
				args:  []interface{}{"world"},
			},
			want: testEncodeWant{
				lineRegexExpr: fmt.Sprintf(
					`^{"datetime":"%s","timestamp":"%s","level":"%s","file":"%s","func":"%s","message":"%s"}\n$`,
					datetimeRegex, timestampRegex, levelRegex, fileCallerRegex, functionCallerRegex, messageRegex,
				),
			},
		},
		{
			args: testEncodeArgs{
				cfg: Config{
					Fields:    []Field{{"foo", "bar"}},
					UTC:       true,
					Datetime:  true,
					Timestamp: true,
					Shortfile: true,
					Longfile:  true,
					Function:  true,
				},
				level: DEBUG,
				msg:   "Hello %s",
				args:  []interface{}{"world"},
			},
			want: testEncodeWant{
				lineRegexExpr: fmt.Sprintf(
					`^{"datetime":"%s","timestamp":"%s","level":"%s","file":"%s","func":"%s",%s,"message":"%s"}\n$`,
					datetimeRegex, timestampRegex, levelRegex, fileCallerRegex,
					functionCallerRegex, fieldsJSONRegex, messageRegex,
				),
			},
		},
		{
			args: testEncodeArgs{
				cfg: Config{
					Fields:    []Field{{"foo", `id: "bar"`}},
					UTC:       true,
					Datetime:  true,
					Timestamp: true,
					Shortfile: true,
					Longfile:  true,
					Function:  true,
				},
				level: DEBUG,
				msg:   "Hello %s",
				args:  []interface{}{"world"},
			},
			want: testEncodeWant{
				lineRegexExpr: fmt.Sprintf(
					`^{"datetime":"%s","timestamp":"%s","level":"%s","file":"%s","func":"%s",%s,"message":"%s"}\n$`,
					datetimeRegex, timestampRegex, levelRegex, fileCallerRegex,
					functionCallerRegex, fieldsJSONRegex, messageRegex,
				),
			},
		},
		{
			args: testEncodeArgs{
				cfg: Config{
					Fields:    []Field{{"foo", "bar"}},
					UTC:       true,
					Datetime:  true,
					Timestamp: true,
					Shortfile: true,
					Longfile:  true,
					Function:  true,
				},
				level: PRINT,
				msg:   "Hello %s",
				args:  []interface{}{"world"},
			},
			want: testEncodeWant{
				lineRegexExpr: fmt.Sprintf(
					`^{"datetime":"%s","timestamp":"%s","file":"%s","func":"%s",%s,"message":"%s"}\n$`,
					datetimeRegex, timestampRegex, fileCallerRegex,
					functionCallerRegex, fieldsJSONRegex, messageRegex,
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
