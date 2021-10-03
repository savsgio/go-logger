package logger

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

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

	outputPtr := reflect.ValueOf(l.output).Pointer()
	wantOutputPtr := reflect.ValueOf(output).Pointer()

	if outputPtr != wantOutputPtr {
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
	levels := []Level{PRINT, FATAL, ERROR, WARNING, DEBUG}

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

func BenchmarkInfo(b *testing.B) {
	l := New(INFO, ioutil.Discard)
	l.SetEncoder(NewEncoderJSON())
	l.SetFlags(Ldatetime | Ltimestamp)
	l.SetFields(Field{Key: "hola", Value: 1}, Field{Key: "adios", Value: 2})

	b.Run("lineal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			l.Info("hello %s", "world")
		}
	})

	b.Run("parallel", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				l.Info("hello %s", "world")
			}
		})
	})
}
