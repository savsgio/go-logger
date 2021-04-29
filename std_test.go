package logger

import (
	"bytes"
	"log"
	"reflect"
	"strings"
	"testing"
)

func TestStd_SetLevel(t *testing.T) {
	level := DEBUG

	SetLevel(level)

	if std.level != debugLevel {
		t.Errorf("Logger.SetLevel() = %v, want %v", std.level, debugLevel)
	}
}

func TestStd_LevelEnabled(t *testing.T) { // nolint:funlen
	type args struct {
		level string
	}

	type want struct {
		fatalEnabled   bool
		errorEnabled   bool
		warningEnabled bool
		infoEnabled    bool
		debugEnabled   bool
	}

	tests := []struct { // nolint:dupl
		name string
		args args
		want want
	}{
		{
			name: "Fatal",
			args: args{level: FATAL},
			want: want{
				fatalEnabled:   true,
				errorEnabled:   false,
				warningEnabled: false,
				infoEnabled:    false,
				debugEnabled:   false,
			},
		},
		{
			name: "Error",
			args: args{level: ERROR},
			want: want{
				fatalEnabled:   true,
				errorEnabled:   true,
				warningEnabled: false,
				infoEnabled:    false,
				debugEnabled:   false,
			},
		},
		{
			name: "Warning",
			args: args{level: WARNING},
			want: want{
				fatalEnabled:   true,
				errorEnabled:   true,
				warningEnabled: true,
				infoEnabled:    false,
				debugEnabled:   false,
			},
		},
		{
			name: "Info",
			args: args{level: INFO},
			want: want{
				fatalEnabled:   true,
				errorEnabled:   true,
				warningEnabled: true,
				infoEnabled:    true,
				debugEnabled:   false,
			},
		},
		{
			name: "Debug",
			args: args{level: DEBUG},
			want: want{
				fatalEnabled:   true,
				errorEnabled:   true,
				warningEnabled: true,
				infoEnabled:    true,
				debugEnabled:   true,
			},
		},
	}
	for _, tt := range tests {
		test := tt

		t.Run(test.name, func(t *testing.T) {
			SetLevel(test.args.level)

			isEnabled := FatalEnabled()
			if isEnabled != test.want.fatalEnabled {
				t.Errorf("FatalEnabled() = '%v', want '%v'", isEnabled, test.want.fatalEnabled)
			}

			isEnabled = ErrorEnabled()
			if isEnabled != test.want.errorEnabled {
				t.Errorf("ErrorEnabled() = '%v', want '%v'", isEnabled, test.want.errorEnabled)
			}

			isEnabled = WarningEnabled()
			if isEnabled != test.want.warningEnabled {
				t.Errorf("WarningEnabled() = '%v', want '%v'", isEnabled, test.want.warningEnabled)
			}

			isEnabled = InfoEnabled()
			if isEnabled != test.want.infoEnabled {
				t.Errorf("InfoEnabled() = '%v', want '%v'", isEnabled, test.want.infoEnabled)
			}

			isEnabled = DebugEnabled()
			if isEnabled != test.want.debugEnabled {
				t.Errorf("DebugEnabled() = '%v', want '%v'", isEnabled, test.want.debugEnabled)
			}
		})
	}
}

func TestStd_SetOutput(t *testing.T) {
	output := new(bytes.Buffer)

	SetOutput(output)

	if !reflect.DeepEqual(std.out, output) {
		t.Errorf("Logger.SetOutput() = %p, want %p", std.out, output)
	}
}

func TestStd_SetFlags(t *testing.T) {
	flags := log.Ldate | log.Ltime | log.Llongfile
	SetFlags(flags)

	if !reflect.DeepEqual(std.flag, flags) {
		t.Errorf("Logger.SetOutput() = %d, want %d", std.flag, flags)
	}
}

func TestStdCallDepth(t *testing.T) {
	output := new(bytes.Buffer)

	SetFlags(log.Lshortfile)
	SetOutput(output)

	Print("Calldepth path test")

	if got := output.String(); !strings.HasPrefix(got, "std.go") {
		t.Errorf("Logger.Print() = %v, want to start with std.go", got)
	}

	output.Reset()
}

func TestStdErrorAndErrorf(t *testing.T) {
	output := new(bytes.Buffer)

	SetLevel(ERROR)
	SetOutput(output)

	t.Run("Error", func(t *testing.T) {
		Error("Error msg")

		if len(output.Bytes()) == 0 {
			t.Error("Error() test failed")
		}
		output.Reset()
	})

	t.Run("Errorf", func(t *testing.T) {
		Errorf("Error msg with %s", "params")

		if len(output.Bytes()) == 0 {
			t.Error("Error() test failed")
		}
		output.Reset()
	})
}

func TestStdWarningAndWarningf(t *testing.T) {
	output := new(bytes.Buffer)

	SetLevel(WARNING)
	SetOutput(output)

	t.Run("Warning", func(t *testing.T) {
		Warning("Warning msg")

		if len(output.Bytes()) == 0 {
			t.Error("Warning() test failed")
		}
		output.Reset()
	})

	t.Run("Warningf", func(t *testing.T) {
		Warningf("Warning msg with %s", "params")

		if len(output.Bytes()) == 0 {
			t.Error("Warningf() test failed")
		}
		output.Reset()
	})
}

func TestStdInfoAndInfof(t *testing.T) {
	output := new(bytes.Buffer)

	SetLevel(INFO)
	SetOutput(output)

	t.Run("Info", func(t *testing.T) {
		Info("Info msg")

		if len(output.Bytes()) == 0 {
			t.Error("Info() test failed")
		}
		output.Reset()
	})

	t.Run("Infof", func(t *testing.T) {
		Infof("Info msg with %s", "params")

		if len(output.Bytes()) == 0 {
			t.Error("Infof() test failed")
		}
		output.Reset()
	})
}

func TestStdDebugAndDebugf(t *testing.T) {
	output := new(bytes.Buffer)

	SetLevel(DEBUG)
	SetOutput(output)

	t.Run("Debug", func(t *testing.T) {
		Debug("Debug msg")

		if len(output.Bytes()) == 0 {
			t.Error("Debug() test failed")
		}
		output.Reset()
	})

	t.Run("Debugf", func(t *testing.T) {
		Debugf("Debug msg with %s", "params")

		if len(output.Bytes()) == 0 {
			t.Error("Debugf() test failed")
		}
		output.Reset()
	})
}

func TestStdPrintAndPrintf(t *testing.T) {
	output := new(bytes.Buffer)

	SetLevel(DEBUG)
	SetOutput(output)

	t.Run("Print", func(t *testing.T) {
		Print("Print msg")
		if len(output.Bytes()) == 0 {
			t.Error("Print() test failed")
		}
		output.Reset()
	})

	t.Run("Printf", func(t *testing.T) {
		Printf("Print msg with %s", "params")
		if len(output.Bytes()) == 0 {
			t.Error("Printf() test failed")
		}
		output.Reset()
	})
}
