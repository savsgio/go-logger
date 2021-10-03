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

	if cfg := l.encoder.Config(); !reflect.DeepEqual(cfg, l.cfg) {
		t.Errorf("Logger.enconder.Config() == %v, want %v", cfg, l.cfg)
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
	calldepth := 123

	l := newTestLogger()
	l.setCalldepth(calldepth)

	if l.cfg.calldepth != calldepth {
		t.Errorf("calldepth == %d, want %d", l.cfg.calldepth, calldepth)
	}

	if cfg := l.encoder.Config(); cfg.calldepth != calldepth {
		t.Errorf("calldepth == %d, want %d", l.cfg.calldepth, calldepth)
	}
}

func TestLogger_setFields(t *testing.T) {
	field1 := Field{"key", "value"}
	field2 := Field{"foo", []int{1, 2, 3}}

	l := newTestLogger()
	l.setFields(field1, field2)

	totalFields := len(l.cfg.Fields)
	wantTotalFields := 2

	if totalFields != wantTotalFields {
		t.Errorf("length == %d, want %d", totalFields, wantTotalFields)
	}

	if cfg := l.encoder.Config(); !reflect.DeepEqual(l.cfg.Fields, cfg.Fields) {
		t.Errorf("result == %v, want %v", l.cfg.Fields, cfg.Fields)
	}

	newField1 := Field{field1.Key, 123.45}
	l.setFields(newField1)

	if currentTotalFields := len(l.cfg.Fields); currentTotalFields != wantTotalFields {
		t.Errorf("length (update) == %d, want %d", currentTotalFields, wantTotalFields)
	}

	if field := l.getField(field1.Key); reflect.DeepEqual(field.Value, field1.Value) {
		t.Errorf("field1 value (update) == %v, want %v", field.Value, newField1.Value)
	}

	newField := Field{"data", []interface{}{1, "2", nil}}
	l.setFields(newField)

	wantTotalFields++

	if currentTotalFields := len(l.cfg.Fields); currentTotalFields != wantTotalFields {
		t.Errorf("length (append) == %d, want %d", currentTotalFields, wantTotalFields)
	}

	if field := l.getField(newField.Key); !reflect.DeepEqual(field.Value, newField.Value) {
		t.Error("the field value is not updated")
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
