package logger

import (
	"bytes"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

var levels = []Level{PRINT, FATAL, ERROR, WARNING, INFO, DEBUG}

func newTestLogger() *Logger {
	return New(DEBUG, ioutil.Discard)
}

func Test_New(t *testing.T) {
	level := INFO
	output := os.Stderr
	fields := []Field{{"key", "value"}}

	l := New(level, output, fields...)

	wantCfg := EncoderConfig{
		Fields:    fields,
		Flag:      LstdFlags,
		Datetime:  true,
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

	if encoderCfg := l.encoder.Config(); !reflect.DeepEqual(encoderCfg, l.cfg) {
		t.Errorf("Logger.enconder.Config() == %v, want %v", encoderCfg, l.cfg)
	}

	if l.encodeOutput == nil {
		t.Fatal("Logger.encodeOutput is nil")
	}
}

func TestLogger_encodeOutput(t *testing.T) {
	output := new(bytes.Buffer)

	l := newTestLogger()
	l.SetOutput(output)
	l.SetLevel(INFO)

	l.encodeOutput(ERROR, infoLevelStr, "hello %s", []interface{}{"word"})

	if output.Len() == 0 {
		t.Error("enconded output has not been written")
	}

	output.Reset()

	l.encodeOutput(DEBUG, debugLevelStr, "hello %s", []interface{}{"word"})

	if output.Len() > 0 {
		t.Error("enconded output has been written")
	}
}

func TestLogger_getField(t *testing.T) {
	field := Field{Key: "key", Value: "value"}

	l := newTestLogger()
	l.setFields(field)

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

	if encoderCfg := l.encoder.Config(); encoderCfg.calldepth != testCalldepth {
		t.Errorf("encoder calldepth == %d, want %d", encoderCfg.calldepth, testCalldepth)
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

			if encoderCfg := l.encoder.Config(); !reflect.DeepEqual(encoderCfg.Fields, l.cfg.Fields) {
				t.Errorf("encoder fields == %v, want %v", encoderCfg.Fields, l.cfg.Fields)
			}

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

func TestLogger_clone(t *testing.T) {
	l1 := newTestLogger()
	l1.SetOutput(new(bytes.Buffer))

	l2 := l1.clone()

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
}

func TestLogger_WithFields(t *testing.T) {
	l1 := newTestLogger()
	l2 := l1.WithFields(Field{"key", "value"})

	l1TotalFields := len(l1.cfg.Fields)
	l2TotalFields := len(l2.cfg.Fields)

	if l2TotalFields == l1TotalFields {
		t.Errorf("fields == %d, want %d", l2TotalFields, l1TotalFields+1)
	}
}

func TestLogger_SetFields(t *testing.T) {
	fields := []Field{{"key", "value"}, {"foo", "bar"}}
	l1 := newTestLogger()

	beforeTotalFields := len(l1.cfg.Fields)

	l1.SetFields(fields...)

	afterTotalFields := len(l1.cfg.Fields)

	if beforeTotalFields == afterTotalFields {
		t.Errorf("fields == %d, want %d", afterTotalFields, beforeTotalFields+len(fields))
	}
}

func TestLogger_SetFlags(t *testing.T) { // nolint:funlen
	type args struct {
		flag Flag
	}

	type want struct {
		cfg EncoderConfig
	}

	l := newTestLogger()

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "datetime",
			args: args{flag: Ldatetime},
			want: want{
				cfg: EncoderConfig{
					Flag:     Ldatetime,
					Datetime: true,
				},
			},
		},
		{
			name: "timestamp",
			args: args{flag: Ltimestamp},
			want: want{
				cfg: EncoderConfig{
					Flag:      Ltimestamp,
					Timestamp: true,
				},
			},
		},
		{
			name: "utc",
			args: args{flag: LUTC},
			want: want{
				cfg: EncoderConfig{
					Flag: LUTC,
					UTC:  true,
				},
			},
		},
		{
			name: "shortfile",
			args: args{flag: Lshortfile},
			want: want{
				cfg: EncoderConfig{
					Flag:      Lshortfile,
					Shortfile: true,
				},
			},
		},
		{
			name: "longfile",
			args: args{flag: Llongfile},
			want: want{
				cfg: EncoderConfig{
					Flag:     Llongfile,
					Longfile: true,
				},
			},
		},
		{
			name: "std",
			args: args{flag: LstdFlags},
			want: want{
				cfg: EncoderConfig{
					Flag:     LstdFlags,
					Datetime: true,
				},
			},
		},
		{
			name: "all",
			args: args{flag: Ldatetime | Ltimestamp | LUTC | Llongfile | Lshortfile},
			want: want{
				cfg: EncoderConfig{
					Flag:      Ldatetime | Ltimestamp | LUTC | Llongfile | Lshortfile,
					Datetime:  true,
					Timestamp: true,
					UTC:       true,
					Shortfile: true,
					Longfile:  true,
				},
			},
		},
	}

	for i := range tests {
		test := tests[i]

		t.Run(test.name, func(t *testing.T) {
			t.Helper()

			l.SetFlags(test.args.flag)

			if l.cfg.Flag != test.want.cfg.Flag {
				t.Errorf("Flag == %d, want %d", l.cfg.Flag, test.want.cfg.Flag)
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

			if enconderCfg := l.encoder.Config(); !reflect.DeepEqual(enconderCfg, l.cfg) {
				t.Errorf("enconder config == %v, want %v", enconderCfg, l.cfg)
			}
		})
	}
}

func TestLogger_SetLevel(t *testing.T) {
	level := DEBUG

	l1 := newTestLogger()
	l1.SetLevel(level)

	if l1.level != level {
		t.Errorf("level == %d, want %d", l1.level, level)
	}
}

func TestLogger_SetOutput(t *testing.T) {
	output := new(bytes.Buffer)

	l1 := newTestLogger()
	l1.SetOutput(output)

	if l1.output != output {
		t.Errorf("output == %p, want %p", l1.output, output)
	}
}

func TestLogger_SetEncoder(t *testing.T) {
	encoder := NewEncoderJSON()

	l1 := newTestLogger()
	l1.SetEncoder(encoder)

	if l1.encoder != encoder {
		t.Errorf("encoder == %p, want %p", l1.encoder, encoder)
	}
}

func TestLogger_IsLevelEnabled(t *testing.T) {
	l1 := newTestLogger()
	l1.SetLevel(ERROR)

	if !l1.IsLevelEnabled(FATAL) {
		t.Error("level is not enabled")
	}

	if !l1.IsLevelEnabled(ERROR) {
		t.Error("level is not enabled")
	}

	if l1.IsLevelEnabled(DEBUG) {
		t.Error("level is enabled")
	}
}

func TestLogger_Levels(t *testing.T) { // nolint:funlen
	type args struct {
		fn  func(msg ...interface{})
		fnf func(msg string, args ...interface{})
	}

	type want struct {
		level    Level
		levelStr string
	}

	type loggerWrapper struct {
		*Logger

		encodeLevel    Level
		encodeLevelStr string
		encodeMsg      string
		encodeArgs     []interface{}
	}

	l := &loggerWrapper{
		Logger: newTestLogger(),
	}
	l.Logger.encodeOutput = func(level Level, levelStr, msg string, args []interface{}) {
		l.encodeLevel = level
		l.encodeLevelStr = levelStr
		l.encodeMsg = msg
		l.encodeArgs = args
	}

	resetLoggerWrapper := func(l *loggerWrapper) {
		l.encodeLevel = invalid
		l.encodeLevelStr = ""
		l.encodeMsg = ""
		l.encodeArgs = nil
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Print",
			args: args{
				fn:  l.Print,
				fnf: l.Printf,
			},
			want: want{
				level:    PRINT,
				levelStr: printLevelStr,
			},
		},
		// {
		// 	name: "Fatal",
		// 	args: args{
		// 		fn:  l.Fatal,
		// 		fnf: l.Fatalf,
		// 	},
		// 	want: want{
		// 		level:    FATAL,
		// 		levelStr: fatalLevelStr,
		// 	},
		// },
		{
			name: "Error",
			args: args{
				fn:  l.Error,
				fnf: l.Errorf,
			},
			want: want{
				level:    ERROR,
				levelStr: errorLevelStr,
			},
		},
		{
			name: "Warning",
			args: args{
				fn:  l.Warning,
				fnf: l.Warningf,
			},
			want: want{
				level:    WARNING,
				levelStr: warningLevelStr,
			},
		},
		{
			name: "Info",
			args: args{
				fn:  l.Info,
				fnf: l.Infof,
			},
			want: want{
				level:    INFO,
				levelStr: infoLevelStr,
			},
		},
		{
			name: "Debug",
			args: args{
				fn:  l.Debug,
				fnf: l.Debugf,
			},
			want: want{
				level:    DEBUG,
				levelStr: debugLevelStr,
			},
		},
	}

	assert := func(msg string, args []interface{}, want want) {
		if l.encodeLevel != want.level {
			t.Errorf("level == %d, want %d", l.encodeLevel, want.level)
		}

		if l.encodeLevelStr != want.levelStr {
			t.Errorf("level string == %s, want %s", l.encodeLevelStr, want.levelStr)
		}

		if l.encodeMsg != msg {
			t.Errorf("msg == %s, want %s", l.encodeMsg, msg)
		}

		if !reflect.DeepEqual(l.encodeArgs, args) {
			t.Errorf("args == %s, want %s", l.encodeArgs, args)
		}
	}

	for i := range tests {
		test := tests[i]

		t.Run(test.name, func(t *testing.T) {
			t.Helper()

			msg := ""
			args := []interface{}{"Hello", "world"}

			test.args.fn(args...)
			assert(msg, args, test.want)
		})

		resetLoggerWrapper(l)

		t.Run(test.name+"f", func(t *testing.T) {
			t.Helper()

			msg := "Hello %s"
			args := []interface{}{"world"}

			test.args.fnf(msg, args...)
			assert(msg, args, test.want)
		})

		resetLoggerWrapper(l)
	}
}

func BenchmarkInfo(b *testing.B) {
	l := New(INFO, ioutil.Discard)
	l.SetEncoder(NewEncoderJSON())
	l.SetFlags(Ldatetime | Ltimestamp)
	l.SetFields(Field{Key: "hola", Value: 1}, Field{Key: "adios", Value: 2})

	b.Run("lineal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			l.Infof("hello %s", "world")
		}
	})

	b.Run("parallel", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				l.Infof("hello %s", "world")
			}
		})
	})
}
