package logger

import (
	"reflect"
	"testing"
)

func newTestConfig() Config {
	return Config{
		Fields: []Field{
			{"url", `GET "https://example.com"`},
		},
		Datetime:  true,
		Timestamp: true,
		UTC:       true,
		Shortfile: true,
		Longfile:  false,
		Function:  true,
		flag:      Ldatetime | Ltimestamp | LUTC | Lshortfile | Lfunction,
		calldepth: calldepth,
	}
}

func testConfigCopy(t *testing.T, cfg1, cfg2 Config) {
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

func TestConfig_Copy(t *testing.T) {
	cfg := newTestConfig()
	copyCfg := cfg.Copy()

	testConfigCopy(t, cfg, copyCfg)
}
