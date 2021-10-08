package logger

import (
	"fmt"
	"reflect"
	"testing"
)

func newTestEncoderJSON() *EncoderJSON {
	return NewEncoderJSON()
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

	if !reflect.DeepEqual(copyEnc.EncoderBase, enc.EncoderBase) {
		t.Errorf("EncoderBase == %v, want %v", copyEnc.EncoderBase, enc.EncoderBase)
	}
}

func TestEncoderJSON_SetConfig(t *testing.T) {
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
				fieldsEncoded: "\"foo\":\"bar\",\"buzz\":\"[1 2 3]\",",
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

func TestEncoderJSON_Encode(t *testing.T) {
	testCases := []testEncodeCase{
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
	}

	enc := newTestEncoderJSON()

	testEncoderEncode(t, enc, testCases)
}
