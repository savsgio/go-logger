package logger

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func newTestEncoderText() *EncoderText {
	cfg := newTestConfig()

	enc := NewEncoderText(EncoderTextConfig{})
	enc.Configure(cfg)

	return enc
}

func Test_NewEncoderText(t *testing.T) {
	type args struct {
		cfg EncoderTextConfig
	}

	type want struct {
		cfg EncoderTextConfig
	}

	tests := []struct {
		args args
		want want
	}{
		{
			args: args{
				cfg: EncoderTextConfig{},
			},
			want: want{
				cfg: EncoderTextConfig{
					Separator:       defaultTextSeparator,
					DatetimeLayout:  defaultDatetimeLayout,
					TimestampFormat: defaultTimestampFormat,
				},
			},
		},
		{
			args: args{
				cfg: EncoderTextConfig{
					Separator:       "#",
					DatetimeLayout:  time.RFC1123,
					TimestampFormat: TimestampFormatNanoseconds,
				},
			},
			want: want{
				cfg: EncoderTextConfig{
					Separator:       "#",
					DatetimeLayout:  time.RFC1123,
					TimestampFormat: TimestampFormatNanoseconds,
				},
			},
		},
	}

	for i := range tests {
		test := tests[i]

		t.Run("", func(t *testing.T) {
			enc := NewEncoderText(test.args.cfg)
			if enc == nil {
				t.Fatal("return nil")
			}

			if !reflect.DeepEqual(enc.cfg, test.want.cfg) {
				t.Errorf("confg == %v, want %v", enc.cfg, test.want.cfg)
			}
		})
	}
}

func TestEncoderText_Copy(t *testing.T) {
	enc := newTestEncoderText()
	copyEnc, ok := enc.Copy().(*EncoderText)

	if !ok {
		t.Fatal("the copy is not a EncoderText pointer")
	}

	encPtr := reflect.ValueOf(enc).Pointer()
	copyEncPtr := reflect.ValueOf(copyEnc).Pointer()

	if copyEncPtr == encPtr {
		t.Error("the copy has the same pointer than original")
	}

	testEncoderBaseCopy(t, &enc.EncoderBase, &copyEnc.EncoderBase)
}

func TestEncoderText_Configure(t *testing.T) {
	type args struct {
		cfg Config
	}

	type want struct {
		fieldsEncoded string
	}

	enc := newTestEncoderText()

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
					Fields: []Field{{"foo", "bar"}, {"buzz", []int{1, 2, 3}}},
				},
			},
			want: want{
				fieldsEncoded: "foo=bar" + enc.cfg.Separator + "buzz=[1 2 3]" + enc.cfg.Separator,
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

func TestEncoderText_Encode(t *testing.T) { // nolint:funlen,dupl
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
					"^%s - %s\n$",
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
					"^%s - %s - %s - %s - %s - %s\n$",
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
					"^%s - %s - %s - %s - %s - %s - %s\n$",
					datetimeRegex, timestampRegex, levelRegex, fileCallerRegex,
					functionCallerRegex, fieldsTextRegex, messageRegex,
				),
			},
		},
		{ // print/printf case
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
					"^%s - %s - %s - %s - %s - %s\n$",
					datetimeRegex, timestampRegex, fileCallerRegex,
					functionCallerRegex, fieldsTextRegex, messageRegex,
				),
			},
		},
	}

	enc := newTestEncoderText()

	testEncoderEncode(t, enc, testCases)
}

func BenchmarkEncoderText_Encode(b *testing.B) {
	enc := newTestEncoderText()
	benchmarkEncoderEncode(b, enc)
}
