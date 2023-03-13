package logger

import (
	"reflect"
	"regexp"
	"runtime"
	"testing"
	"time"
)

const (
	datetimeRegex       = `([2-9](\d{3})-(\d{2})-(\d{2})T(\d{2}):(\d{2}):(\d{2})Z)`
	timestampRegex      = `(\d+)`
	levelRegex          = `([A-Z]+)`
	fileCallerRegex     = `((.*)\.go\:\d+)`
	functionCallerRegex = `([a-zA-Z0-9\._\-\/]+)`
	fieldsTextRegex     = `(((.*)=(.*))+)`
	fieldsJSONRegex     = `((\"(fields\.)?(.*)\"\:\"(.*)\")+)`
	messageRegex        = `(.*)`
)

type testEncodeArgs struct {
	cfg   Config
	level Level
	msg   string
	args  []interface{}
}

type testEncodeWant struct {
	lineRegexExpr string
}

type testEncodeCase struct {
	args testEncodeArgs
	want testEncodeWant
}

type mockEncoder struct {
	copy             func() Encoder
	fieldsEncoded    func() string
	setFieldsEncoded func(string)
	configure        func(Config)
	encode           func(*Buffer, Entry) error
}

func (enc *mockEncoder) Copy() Encoder {
	return enc.copy()
}

func (enc *mockEncoder) FieldsEncoded() string {
	return enc.fieldsEncoded()
}

func (enc *mockEncoder) SetFieldsEncoded(fieldsEncoded string) {
	enc.setFieldsEncoded(fieldsEncoded)
}

func (enc *mockEncoder) Configure(cfg Config) {
	enc.configure(cfg)
}

func (enc *mockEncoder) Encode(buf *Buffer, e Entry) error {
	return enc.encode(buf, e)
}

func testEncoderEncode(t *testing.T, enc Encoder, testCases []testEncodeCase) {
	t.Helper()

	for i := range testCases {
		test := testCases[i]
		cfg := test.args.cfg

		now := time.Now()
		if cfg.UTC {
			now = now.UTC()
		}

		var caller runtime.Frame
		if cfg.Shortfile || cfg.Longfile {
			caller = getFileCaller(3)
		}

		t.Run("", func(t *testing.T) {
			t.Helper()

			buf := AcquireBuffer()
			defer ReleaseBuffer(buf)

			enc.Configure(cfg)

			e := Entry{
				Config:  cfg,
				Time:    now,
				Level:   test.args.level,
				Caller:  caller,
				Message: buf.formatMessage(test.args.msg, test.args.args),
			}

			if err := enc.Encode(buf, e); err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			re := regexp.MustCompile(test.want.lineRegexExpr)

			if line := buf.String(); !re.MatchString(line) {
				t.Errorf("line == %s, want regex exp %s", line, test.want.lineRegexExpr)
			}

			buf.Reset()
		})
	}
}

func benchmarkEncoderEncode(b *testing.B, enc Encoder) {
	b.Helper()

	buf := AcquireBuffer()
	defer ReleaseBuffer(buf)

	e := Entry{
		Config:  newTestConfig(),
		Time:    time.Now().UTC(),
		Level:   DEBUG,
		Caller:  getFileCaller(4),
		Message: `failed to request: jojoj""""`,
	}

	enc.Configure(e.Config)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := enc.Encode(buf, e); err != nil {
			b.Fatal(err)
		}

		buf.Reset()
	}
}

func newTestEncoderBase() *EncoderBase {
	enc := new(EncoderBase)

	return enc
}

func Test_newEncoderBase(t *testing.T) {
	if enc := newEncoderBase(); enc == nil {
		t.Error("return nil")
	}
}

func testEncoderBaseCopy(t *testing.T, enc1, enc2 *EncoderBase) {
	t.Helper()

	enc1Ptr := reflect.ValueOf(enc1).Pointer()
	enc2Ptr := reflect.ValueOf(enc2).Pointer()

	if enc1Ptr == enc2Ptr {
		t.Error("the copy has the same pointer than original")
	}

	if enc1.fieldsEncoded != enc2.fieldsEncoded {
		t.Errorf("fieldsEncoded == %v, want %v", enc1.fieldsEncoded, enc2.fieldsEncoded)
	}
}

func TestEncoderBase_Copy(t *testing.T) {
	enc := newTestEncoderBase()
	copyEnc := enc.Copy()

	testEncoderBaseCopy(t, enc, copyEnc)
}

func TestEncoderBase_FieldsEncoded(t *testing.T) {
	enc := newTestEncoderBase()
	enc.fieldsEncoded = "v1 - v2 - v3"

	if fieldsEncoded := enc.FieldsEncoded(); enc.fieldsEncoded != fieldsEncoded {
		t.Errorf("fieldsEncoded == %s, want %s", enc.fieldsEncoded, fieldsEncoded)
	}
}

func TestEncoderBase_SetFieldsEncoded(t *testing.T) {
	fieldsEncoded := "v1 - v2 - v3"

	enc := newTestEncoderBase()
	enc.SetFieldsEncoded(fieldsEncoded)

	if enc.fieldsEncoded != fieldsEncoded {
		t.Errorf("fieldsEncoded == %s, want %s", enc.fieldsEncoded, fieldsEncoded)
	}
}
