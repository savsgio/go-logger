package logger

import (
	"fmt"
	"reflect"
	"testing"
)

func newTestEncoderJSON() *EncoderJSON {
	cfg := newTestConfig()

	enc := NewEncoderJSON()
	enc.SetFields(cfg.Fields)

	return enc
}

func Test_NewEncoderJSON(t *testing.T) {
	if enc := NewEncoderJSON(); enc == nil {
		t.Error("return nil")
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

func TestEncoderJSON_SetFields(t *testing.T) { // nolint:funlen
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
				fieldsEncoded: `"foo":"bar","buzz":"[1 2 3]",`,
			},
		},
		{
			args: args{
				fields: []Field{{"foo", `id: "123"`}, {"buzz", []int{1, 2, 3}}},
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
			enc.SetFields(test.args.fields)

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
				},
				level: DEBUG,
				msg:   "Hello %s",
				args:  []interface{}{"world"},
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
					`^{"datetime":"%s","timestamp":"%s","level":"%s","file":"%s",%s,"message":"%s"}\n$`,
					datetimeRegex, timestampRegex, levelRegex, fileCallerRegex, fieldsKVRegex, messageRegex,
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
				},
				level: DEBUG,
				msg:   "Hello %s",
				args:  []interface{}{"world"},
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
