package logger

import (
	"reflect"
	"testing"
)

func newTestEncoderConfig() EncoderConfig {
	return EncoderConfig{
		Fields: []Field{
			{"test", true},
		},
	}
}

func testEncoderConfigCopy(t *testing.T, cfg1, cfg2 EncoderConfig) {
	t.Helper()

	cfg1Fields := cfg1.Fields
	cfg2Fields := cfg2.Fields

	cfg1.Fields = nil
	cfg2.Fields = nil

	if !reflect.DeepEqual(cfg1, cfg2) {
		t.Errorf("cfg == %v, want %v", cfg1, cfg2)
	}

	if reflect.ValueOf(cfg1Fields).Pointer() == reflect.ValueOf(cfg2Fields).Pointer() {
		t.Errorf("cfg.Fields has the same pointer")
	}

	if !reflect.DeepEqual(cfg1Fields, cfg2Fields) {
		t.Errorf("cfg.Fields are not equals")
	}
}

func TestEncoderConfig_Copy(t *testing.T) {
	cfg := newTestEncoderConfig()
	copyCfg := cfg.Copy()

	testEncoderConfigCopy(t, cfg, copyCfg)
}
