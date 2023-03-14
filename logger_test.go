package logger

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"testing"
)

var levels = []Level{PRINT, TRACE, FATAL, ERROR, WARNING, INFO, DEBUG}

type testLoggerLevelArgs struct {
	fn  func(msg ...interface{})
	fnf func(msg string, args ...interface{})
}

type testLoggerLevelWant struct {
	level    Level
	exitCode int
}

type testLoggerLevelCase struct {
	name string
	args testLoggerLevelArgs
	want testLoggerLevelWant
}

func newTestLogger() *Logger {
	cfg := newTestConfig()

	l := New(DEBUG, ioutil.Discard, cfg.Fields...)
	l.setCalldepth(cfg.calldepth)
	l.SetFlags(cfg.flag)

	return l
}

func assertEncoder(t *testing.T, cfg Config, enc Encoder) {
	t.Helper()

	if fieldsEncoded := enc.FieldsEncoded(); len(cfg.Fields) > 0 && fieldsEncoded == "" {
		t.Error("Logger.encoder has not encoded fields")
	}
}

func Test_New(t *testing.T) {
	level := INFO
	output := os.Stderr
	fields := []Field{{"key", "value"}}

	l := New(level, output, fields...)

	wantCfg := Config{
		Fields:    fields,
		Datetime:  true,
		flag:      LstdFlags,
		calldepth: calldepth,
	}

	if !reflect.DeepEqual(l.cfg, wantCfg) {
		t.Errorf("Logger.cfg == %v, want %v", l.cfg, wantCfg)
	}

	if l.level != level {
		t.Errorf("Logger.level == %d, want %d", l.level, level)
	}

	if l.output != output {
		t.Errorf("Logger.output == %p, want %p", l.output, output)
	}

	if l.encoder == nil {
		t.Fatal("Logger.enconder is nil")
	}

	if _, ok := l.encoder.(*EncoderText); !ok {
		t.Error("Logger.enconder is not a EncoderText pointer")
	}

	assertEncoder(t, l.cfg, l.encoder)

	loggerExitPtr := reflect.ValueOf(l.exit).Pointer()
	osExitPtr := reflect.ValueOf(os.Exit).Pointer()

	if loggerExitPtr != osExitPtr {
		t.Errorf("exit == %p, want %p", l.exit, os.Exit)
	}
}

func TestLogger_encodeOutput(t *testing.T) { // nolint:funlen
	msg := "hello %s"
	args := []interface{}{"men"}
	level := DEBUG
	calldepth := calldepth - 1
	output := new(bytes.Buffer)

	l := newTestLogger()
	l.SetOutput(output)
	l.setCalldepth(calldepth)

	var wantResult string

	enc := new(mockEncoder)
	enc.configure = func(cfg Config) {}
	enc.encode = func(buf *Buffer, e Entry) error {
		t.Helper()

		if buf == nil {
			t.Error("nil buffer")
		}

		if !reflect.DeepEqual(e.Config, l.cfg) {
			t.Errorf("entry config == %v, want %v", e.Config, l.cfg)
		}

		if e.Time.IsZero() {
			t.Error("entry time is zeo")
		}

		if !e.Time.Equal(e.Time.UTC()) {
			t.Error("entry time is not in UTC")
		}

		if e.Level != level {
			t.Errorf("entry level == %s, want %s", e.Level, level)
		}

		wantFile := "logger_test.go"
		if _, file := filepath.Split(e.Caller.File); file != wantFile {
			t.Errorf("entry caller file == %s, want %s", file, wantFile)
		}

		if e.Caller.Line == 0 {
			t.Errorf("entry caller line is zero")
		}

		wantFunction := "github.com/savsgio/go-logger/v4.TestLogger_encodeOutput"
		if e.Caller.Function != wantFunction {
			t.Errorf("entry caller function == %s, want %s", e.Caller.Function, wantFunction)
		}

		wantResult = buf.formatMessage(msg, args)
		if e.Message != wantResult {
			t.Errorf("entry message == %s, want %s", e.Message, wantResult)
		}

		if e.RawMessage != msg {
			t.Errorf("entry raw message == %s, want %s", e.RawMessage, msg)
		}

		if !reflect.DeepEqual(e.Args, args) {
			t.Errorf("entry args == %v, want %v", e.Args, args)
		}

		_, err := buf.WriteString(e.Message)

		return err
	}

	l.SetEncoder(enc)

	hookFired := false
	hook := &testHook{
		levels: []Level{level},
		fireFunc: func(e Entry) error {
			hookFired = true

			return nil
		},
	}

	if err := l.AddHook(hook); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	l.encodeOutput(level, msg, args)

	if result := output.String(); result != wantResult {
		t.Errorf("output result == %s, want %s", result, wantResult)
	}

	if !hookFired {
		t.Errorf("hook not fired")
	}

	hookFired = false

	output.Reset()
	l.SetLevel(ERROR)

	l.encodeOutput(DEBUG, "hello %s", []interface{}{"word"})

	if output.Len() > 0 {
		t.Error("enconded output has been written")
	}

	if hookFired {
		t.Errorf("hook fired")
	}
}

func TestLogger_getField(t *testing.T) {
	field := Field{Key: "key", Value: "value"}

	l := newTestLogger()
	l.SetFields(field)

	result := l.getField(field.Key)

	if result == nil {
		t.Fatal("nil field")
	}

	if !reflect.DeepEqual(result, &field) {
		t.Errorf("result == %v, want %v", result, &field)
	}
}

func TestLogger_setCalldepth(t *testing.T) {
	testCalldepth := 123

	l := newTestLogger()
	l.setCalldepth(testCalldepth)

	if l.cfg.calldepth != testCalldepth {
		t.Errorf("calldepth == %d, want %d", l.cfg.calldepth, testCalldepth)
	}
}

func TestLogger_setFields(t *testing.T) { // nolint:funlen
	field1 := Field{"key", "value"}
	field2 := Field{"foo", []int{1, 2, 3}}

	type args struct {
		fields []Field
	}

	type want struct {
		totalFields int
	}

	l := newTestLogger()
	l.cfg.Fields = nil

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "new",
			args: args{
				fields: []Field{field1, field2},
			},
			want: want{
				totalFields: 2,
			},
		},
		{
			name: "update",
			args: args{
				fields: []Field{{field1.Key, 123.45}},
			},
			want: want{
				totalFields: 2,
			},
		},
		{
			name: "append",
			args: args{
				fields: []Field{{"data", []interface{}{1, "2", nil}}},
			},
			want: want{
				totalFields: 3,
			},
		},
	}

	for i := range tests {
		test := tests[i]

		t.Run(test.name, func(t *testing.T) {
			t.Helper()

			l.setFields(test.args.fields...)

			if totalFields := len(l.cfg.Fields); totalFields != test.want.totalFields {
				t.Errorf("length == %d, want %d", totalFields, test.want.totalFields)
			}

			assertEncoder(t, l.cfg, l.encoder)

			for _, argField := range test.args.fields {
				if field := l.getField(argField.Key); !reflect.DeepEqual(field.Value, argField.Value) {
					t.Errorf("field value == %v, want %v", field.Value, argField.Value)
				}
			}
		})
	}
}

func TestLogger_isLevelEnabled(t *testing.T) {
	l := newTestLogger()

	for _, level := range levels {
		l.SetLevel(level)

		for _, currentLevel := range levels {
			enabled := l.isLevelEnabled(currentLevel)
			wantEnabled := level >= currentLevel

			if enabled != wantEnabled {
				t.Errorf("enabled (level: %d, current: %d) == %t, want %t", level, currentLevel, enabled, wantEnabled)
			}
		}
	}
}

func TestLogger_copy(t *testing.T) {
	l1 := newTestLogger()
	l1.SetOutput(new(bytes.Buffer))

	l2 := l1.copy()

	l1Fields := l1.cfg.Fields
	l2Fields := l2.cfg.Fields

	l1.cfg.Fields = nil
	l2.cfg.Fields = nil

	if !reflect.DeepEqual(l2.cfg, l1.cfg) {
		t.Errorf("cfg == %v, want %v", l2.cfg, l1.cfg)
	}

	if reflect.ValueOf(l2Fields).Pointer() == reflect.ValueOf(l1Fields).Pointer() {
		t.Errorf("fields values has the same pointer")
	}

	if l2.level != l1.level {
		t.Errorf("level == %d, want %d", l2.level, l1.level)
	}

	if l2.output != l1.output {
		t.Errorf("output == %p, want %p", l2.output, l1.output)
	}

	if l2.encoder == l1.encoder {
		t.Error("encoder values has the same pointer")
	}

	l1EncodeOutputPtr := reflect.ValueOf(l1.encodeOutput).Pointer()
	l2EncodeOutputPtr := reflect.ValueOf(l2.encodeOutput).Pointer()

	if l2EncodeOutputPtr != l1EncodeOutputPtr {
		t.Errorf("encodeOutput == %p, want %p", l2.encodeOutput, l1.encodeOutput)
	}

	l1HooksPtr := reflect.ValueOf(l1.hooks).Pointer()
	l2HooksPtr := reflect.ValueOf(l2.hooks).Pointer()

	if l1HooksPtr == l2HooksPtr {
		t.Error("hooks has the same pointer")
	}

	l1ExitPtr := reflect.ValueOf(l1.exit).Pointer()
	l2ExitPtr := reflect.ValueOf(l2.exit).Pointer()

	if l2ExitPtr != l1ExitPtr {
		t.Errorf("exit == %p, want %p", l2.exit, l1.exit)
	}
}

func testLoggerWithFields(t *testing.T, l1 *Logger, withFieldsFunc func(fields ...Field) *Logger) {
	t.Helper()

	l2 := withFieldsFunc(Field{"key", "value"})

	l1TotalFields := len(l1.cfg.Fields)
	l2TotalFields := len(l2.cfg.Fields)

	if l2TotalFields == l1TotalFields {
		t.Errorf("fields == %d, want %d", l2TotalFields, l1TotalFields+1)
	}
}

func TestLogger_WithFields(t *testing.T) {
	l1 := newTestLogger()
	testLoggerWithFields(t, l1, l1.WithFields)
}

func testLoggerSetFields(t *testing.T, l *Logger, setFieldsFunc func(fields ...Field)) {
	t.Helper()

	fields := []Field{{"key", "value"}, {"foo", "bar"}}

	beforeTotalFields := len(l.cfg.Fields)

	setFieldsFunc(fields...)

	afterTotalFields := len(l.cfg.Fields)

	if beforeTotalFields == afterTotalFields {
		t.Errorf("fields == %d, want %d", afterTotalFields, beforeTotalFields+len(fields))
	}
}

func TestLogger_SetFields(t *testing.T) {
	l := newTestLogger()
	testLoggerSetFields(t, l, l.SetFields)
}

func testLoggerSetFlags(t *testing.T, l *Logger, setFlagsFunc func(flag Flag)) { // nolint:funlen
	t.Helper()

	type args struct {
		flag Flag
	}

	type want struct {
		cfg Config
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "datetime",
			args: args{flag: Ldatetime},
			want: want{
				cfg: Config{
					Datetime: true,
					flag:     Ldatetime,
				},
			},
		},
		{
			name: "timestamp",
			args: args{flag: Ltimestamp},
			want: want{
				cfg: Config{
					Timestamp: true,
					flag:      Ltimestamp,
				},
			},
		},
		{
			name: "utc",
			args: args{flag: LUTC},
			want: want{
				cfg: Config{
					UTC:  true,
					flag: LUTC,
				},
			},
		},
		{
			name: "shortfile",
			args: args{flag: Lshortfile},
			want: want{
				cfg: Config{
					Shortfile: true,
					flag:      Lshortfile,
				},
			},
		},
		{
			name: "longfile",
			args: args{flag: Llongfile},
			want: want{
				cfg: Config{
					Longfile: true,
					flag:     Llongfile,
				},
			},
		},
		{
			name: "function",
			args: args{flag: Lfunction},
			want: want{
				cfg: Config{
					Function: true,
					flag:     Lfunction,
				},
			},
		},
		{
			name: "std",
			args: args{flag: LstdFlags},
			want: want{
				cfg: Config{
					Datetime: true,
					flag:     LstdFlags,
				},
			},
		},
		{
			name: "all",
			args: args{flag: Ldatetime | Ltimestamp | LUTC | Llongfile | Lshortfile | Lfunction},
			want: want{
				cfg: Config{
					Datetime:  true,
					Timestamp: true,
					UTC:       true,
					Shortfile: true,
					Longfile:  true,
					Function:  true,
					flag:      Ldatetime | Ltimestamp | LUTC | Llongfile | Lshortfile | Lfunction,
				},
			},
		},
	}

	for i := range tests {
		test := tests[i]

		t.Run(test.name, func(t *testing.T) {
			t.Helper()

			setFlagsFunc(test.args.flag)

			if l.cfg.flag != test.want.cfg.flag {
				t.Errorf("Flag == %d, want %d", l.cfg.flag, test.want.cfg.flag)
			}

			if l.cfg.Datetime != test.want.cfg.Datetime {
				t.Errorf("Datetime == %t, want %t", l.cfg.Datetime, test.want.cfg.Datetime)
			}

			if l.cfg.Timestamp != test.want.cfg.Timestamp {
				t.Errorf("Timestamp == %t, want %t", l.cfg.Timestamp, test.want.cfg.Timestamp)
			}

			if l.cfg.UTC != test.want.cfg.UTC {
				t.Errorf("UTC == %t, want %t", l.cfg.UTC, test.want.cfg.UTC)
			}

			if l.cfg.Shortfile != test.want.cfg.Shortfile {
				t.Errorf("Shortfile == %t, want %t", l.cfg.Shortfile, test.want.cfg.Shortfile)
			}

			if l.cfg.Longfile != test.want.cfg.Longfile {
				t.Errorf("Longfile == %t, want %t", l.cfg.Longfile, test.want.cfg.Longfile)
			}
		})
	}
}

func TestLogger_SetFlags(t *testing.T) {
	l := newTestLogger()
	testLoggerSetFlags(t, l, l.SetFlags)
}

func testLoggerSetLevel(t *testing.T, l *Logger, setLevelFunc func(level Level)) {
	t.Helper()

	level := DEBUG

	setLevelFunc(level)

	if l.level != level {
		t.Errorf("level == %d, want %d", l.level, level)
	}
}

func TestLogger_SetLevel(t *testing.T) {
	l := newTestLogger()
	testLoggerSetLevel(t, l, l.SetLevel)
}

func testLoggerSetOutput(t *testing.T, l *Logger, setOutputFunc func(output io.Writer)) {
	t.Helper()

	output := new(bytes.Buffer)

	setOutputFunc(output)

	if l.output != output {
		t.Errorf("output == %p, want %p", l.output, output)
	}
}

func TestLogger_SetOutput(t *testing.T) {
	l := newTestLogger()
	testLoggerSetOutput(t, l, l.SetOutput)
}

func testLoggerSetEncoder(t *testing.T, l *Logger, setEncoderFunc func(enc Encoder)) {
	t.Helper()

	encoder := newTestEncoderJSON()

	setEncoderFunc(encoder)

	if l.encoder != encoder {
		t.Errorf("encoder == %p, want %p", l.encoder, encoder)
	}

	assertEncoder(t, l.cfg, l.encoder)
}

func TestLogger_SetEncoder(t *testing.T) {
	l := newTestLogger()
	testLoggerSetEncoder(t, l, l.SetEncoder)
}

func testLoggerIsLevelEnabled(t *testing.T, l *Logger, isLevelEnabledFunc func(level Level) bool) {
	t.Helper()

	l.SetLevel(ERROR)

	if !isLevelEnabledFunc(FATAL) {
		t.Error("level is not enabled")
	}

	if !isLevelEnabledFunc(ERROR) {
		t.Error("level is not enabled")
	}

	if isLevelEnabledFunc(DEBUG) {
		t.Error("level is enabled")
	}
}

func TestLogger_IsLevelEnabled(t *testing.T) {
	l := newTestLogger()
	testLoggerIsLevelEnabled(t, l, l.IsLevelEnabled)
}

func testLoggerAddHook(t *testing.T, l *Logger, addHookFunc func(h Hook) error) {
	t.Helper()

	type args struct {
		hook *testHook
	}

	type want struct {
		err error
	}

	tests := []struct {
		args args
		want want
	}{
		{
			args: args{
				hook: &testHook{
					levels:   levels,
					fireFunc: func(e Entry) error { return nil },
				},
			},
			want: want{
				err: nil,
			},
		},
		{
			args: args{
				hook: &testHook{
					levels: []Level{},
				},
			},
			want: want{
				err: ErrEmptyHookLevels,
			},
		},
	}

	for i := range tests {
		test := tests[i]

		t.Run("", func(t *testing.T) {
			if err := addHookFunc(test.args.hook); !errors.Is(err, test.want.err) {
				t.Errorf("error == %v, want %v", err, test.want.err)
			}

			errorExpected := test.want.err != nil

			if !errorExpected && len(l.hooks.store) == 0 {
				t.Errorf("hook not added")
			}
		})
	}
}

func TestLogger_AddHook(t *testing.T) {
	l := newTestLogger()
	testLoggerAddHook(t, l, l.AddHook)
}

func testLoggerLevels(t *testing.T, l *Logger, testCases []testLoggerLevelCase) { // nolint:funlen
	t.Helper()

	var (
		exitCode = -1
		entry    = Entry{}
	)

	l.exit = func(code int) {
		exitCode = code
	}

	enc := new(mockEncoder)
	enc.configure = func(c Config) {}
	enc.encode = func(b *Buffer, e Entry) error {
		entry = e

		return nil
	}

	l.SetEncoder(enc)

	assert := func(msg string, args []interface{}, want testLoggerLevelWant) {
		if entry.Level != want.level {
			t.Errorf("level == %d, want %d", entry.Level, want.level)
		}

		wantFile := "logger_test.go"
		if _, file := filepath.Split(entry.Caller.File); file != wantFile {
			t.Errorf("entry caller file == %s, want %s", file, wantFile)
		}

		if entry.Caller.Line == 0 {
			t.Errorf("entry caller line is zero")
		}

		wantFunction := regexp.MustCompile("^github.com/savsgio/go-logger/v4.testLoggerLevels.func([0-9]{1})$")
		if !wantFunction.MatchString(entry.Caller.Function) {
			t.Errorf("entry caller function == %s, want %s", entry.Caller.Function, wantFunction)
		}

		if entry.RawMessage != msg {
			t.Errorf("msg == %s, want %s", entry.RawMessage, msg)
		}

		if !reflect.DeepEqual(entry.Args, args) {
			t.Errorf("args == %s, want %s", entry.Args, args)
		}

		if exitCode != want.exitCode {
			t.Errorf("exit code == %d, want %d", exitCode, want.exitCode)
		}

		// reset
		exitCode = -1
		entry = Entry{}
	}

	for i := range testCases {
		test := testCases[i]

		t.Run(test.name, func(t *testing.T) {
			t.Helper()

			msg := ""
			args := []interface{}{"Hello", "world"}

			test.args.fn(args...)
			assert(msg, args, test.want)
		})

		t.Run(test.name+"f", func(t *testing.T) {
			t.Helper()

			msg := "Hello %s"
			args := []interface{}{"world"}

			test.args.fnf(msg, args...)
			assert(msg, args, test.want)
		})
	}
}

func TestLogger_Levels(t *testing.T) { // nolint:funlen
	l := newTestLogger()

	testCases := []testLoggerLevelCase{
		{
			name: "Print",
			args: testLoggerLevelArgs{
				fn:  l.Print,
				fnf: l.Printf,
			},
			want: testLoggerLevelWant{
				level:    PRINT,
				exitCode: -1,
			},
		},
		{
			name: "Trace",
			args: testLoggerLevelArgs{
				fn:  l.Trace,
				fnf: l.Tracef,
			},
			want: testLoggerLevelWant{
				level:    TRACE,
				exitCode: -1,
			},
		},
		{
			name: "Fatal",
			args: testLoggerLevelArgs{
				fn:  l.Fatal,
				fnf: l.Fatalf,
			},
			want: testLoggerLevelWant{
				level:    FATAL,
				exitCode: 1,
			},
		},
		{
			name: "Error",
			args: testLoggerLevelArgs{
				fn:  l.Error,
				fnf: l.Errorf,
			},
			want: testLoggerLevelWant{
				level:    ERROR,
				exitCode: -1,
			},
		},
		{
			name: "Warning",
			args: testLoggerLevelArgs{
				fn:  l.Warning,
				fnf: l.Warningf,
			},
			want: testLoggerLevelWant{
				level:    WARNING,
				exitCode: -1,
			},
		},
		{
			name: "Info",
			args: testLoggerLevelArgs{
				fn:  l.Info,
				fnf: l.Infof,
			},
			want: testLoggerLevelWant{
				level:    INFO,
				exitCode: -1,
			},
		},
		{
			name: "Debug",
			args: testLoggerLevelArgs{
				fn:  l.Debug,
				fnf: l.Debugf,
			},
			want: testLoggerLevelWant{
				level:    DEBUG,
				exitCode: -1,
			},
		},
	}

	testLoggerLevels(t, l, testCases)
}

func BenchmarkLogger_Levels(b *testing.B) { // nolint:funlen
	l := newTestLogger()
	l.SetEncoder(newTestEncoderJSON())
	// l.SetFlags(Ltimestamp)
	l.SetFields(Field{Key: "hola", Value: 1}, Field{Key: "adios", Value: 2})
	l.SetLevel(DEBUG)
	// l.SetFlags(0)

	l.exit = func(code int) {}

	benchs := []struct {
		name string
		args testLoggerLevelArgs
	}{
		{
			name: "Print",
			args: testLoggerLevelArgs{
				fn:  l.Print,
				fnf: l.Printf,
			},
		},
		{
			name: "Trace",
			args: testLoggerLevelArgs{
				fn:  l.Trace,
				fnf: l.Tracef,
			},
		},
		{
			name: "Fatal",
			args: testLoggerLevelArgs{
				fn:  l.Fatal,
				fnf: l.Fatalf,
			},
		},
		{
			name: "Error",
			args: testLoggerLevelArgs{
				fn:  l.Error,
				fnf: l.Errorf,
			},
		},
		{
			name: "Warning",
			args: testLoggerLevelArgs{
				fn:  l.Warning,
				fnf: l.Warningf,
			},
		},
		{
			name: "Info",
			args: testLoggerLevelArgs{
				fn:  l.Info,
				fnf: l.Infof,
			},
		},
		{
			name: "Debug",
			args: testLoggerLevelArgs{
				fn:  l.Debug,
				fnf: l.Debugf,
			},
		},
	}

	b.Run("lineal", func(b *testing.B) {
		for i := range benchs {
			bench := benchs[i]

			b.Run(bench.name, func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					bench.args.fn("hello world")
				}
			})

			b.Run(bench.name+"f", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					bench.args.fnf("hello %s", " world")
				}
			})
		}
	})

	b.Run("parallel", func(b *testing.B) {
		for i := range benchs {
			bench := benchs[i]

			b.Run(bench.name, func(b *testing.B) {
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						bench.args.fn("hello world")
					}
				})
			})

			b.Run(bench.name+"f", func(b *testing.B) {
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						bench.args.fnf("hello %s", " world")
					}
				})
			})
		}
	})
}
