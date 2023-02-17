package logger

import (
	"fmt"
	"reflect"
	"testing"
)

func newTestEncoderText() *EncoderText {
	cfg := newTestConfig()

	enc := NewEncoderText()
	enc.SetFields(cfg.Fields)

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

func TestEncoderText_SetFields(t *testing.T) {
	type args struct {
		fields []Field
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
				fields: []Field{},
			},
			want: want{
				fieldsEncoded: "",
			},
		},
		{
			args: args{
				fields: []Field{{"foo", "bar"}, {"buzz", []int{1, 2, 3}}},
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
			enc.SetFields(test.args.fields)

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
				},
				level: DEBUG,
				msg:   "Hello %s",
				args:  []interface{}{"world"},
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
				cfg: Config{
					Fields:    []Field{{"foo", "bar"}},
					UTC:       true,
					Datetime:  true,
					Timestamp: true,
					Shortfile: true,
					Longfile:  true,
				},
				level: DEBUG,
				msg:   "Hello %s",
				args:  []interface{}{"world"},
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
				cfg: Config{
					Fields:    []Field{{"foo", "bar"}},
					UTC:       true,
					Datetime:  true,
					Timestamp: true,
					Shortfile: true,
					Longfile:  true,
				},
				level: PRINT,
				msg:   "Hello %s",
				args:  []interface{}{"world"},
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
