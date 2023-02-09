package logger

import (
	"fmt"
	"reflect"
	"testing"
)

func newTestEncoderText() *EncoderText {
	enc := NewEncoderText()
	enc.SetConfig(newTestEncoderConfig())

	return enc
}

func Test_NewEncoderText(t *testing.T) {
	if enc := NewEncoderText(); enc == nil {
		t.Error("return nil")
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

func TestEncoderText_SetConfig(t *testing.T) {
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
				fieldsEncoded: "{\"foo\":\"bar\",\"buzz\":\"[1 2 3]\"}" + sepText,
			},
		},
	}

	for i := range tests {
		test := tests[i]

		t.Run("", func(t *testing.T) {
			t.Helper()

			enc := newTestEncoderText()
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

func TestEncoderText_Encode(t *testing.T) { // nolint:funlen,dupl
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
					"^%s - %s\n$",
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
					"^%s - %s - %s - %s - %s\n$",
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
					"^%s - %s - %s - %s - %s - %s\n$",
					datetimeRegex, timestampRegex, levelRegex, fileCallerRegex, fieldsJSONRegex, messageRegex,
				),
			},
		},
		{ // print/printf case
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
					"^%s - %s - %s - %s - %s\n$",
					datetimeRegex, timestampRegex, fileCallerRegex, fieldsJSONRegex, messageRegex,
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
