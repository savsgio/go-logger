package logger

import (
	"testing"
)

func TestLogger_std_WithFields(t *testing.T) {
	testLoggerWithFields(t, std, WithFields)
}

func TestLogger_std_SetFields(t *testing.T) {
	testLoggerSetFields(t, std, SetFields)
}

func TestLogger_std_SetFlags(t *testing.T) {
	testLoggerSetFlags(t, std, SetFlags)
}

func TestLogger_std_SetLevel(t *testing.T) {
	testLoggerSetLevel(t, std, SetLevel)
}

func TestLogger_std_SetOutput(t *testing.T) {
	testLoggerSetOutput(t, std, SetOutput)
}

func TestLogger_std_SetEncoder(t *testing.T) {
	testLoggerSetEncoder(t, std, SetEncoder)
}

func TestLogger_std_IsLevelEnabled(t *testing.T) {
	testLoggerIsLevelEnabled(t, std, IsLevelEnabled)
}

func TestLogger_std_Levels(t *testing.T) { // nolint:funlen
	testCases := []testLoggerLevelCase{
		{
			name: "Print",
			args: testLoggerLevelArgs{
				fn:  Print,
				fnf: Printf,
			},
			want: testLoggerLevelWant{
				level:    PRINT,
				levelStr: printLevelStr,
				exitCode: -1,
			},
		},
		{
			name: "Trace",
			args: testLoggerLevelArgs{
				fn:  Trace,
				fnf: Tracef,
			},
			want: testLoggerLevelWant{
				level:    TRACE,
				levelStr: traceLevelStr,
				exitCode: -1,
			},
		},
		{
			name: "Fatal",
			args: testLoggerLevelArgs{
				fn:  Fatal,
				fnf: Fatalf,
			},
			want: testLoggerLevelWant{
				level:    FATAL,
				levelStr: fatalLevelStr,
				exitCode: 1,
			},
		},
		{
			name: "Error",
			args: testLoggerLevelArgs{
				fn:  Error,
				fnf: Errorf,
			},
			want: testLoggerLevelWant{
				level:    ERROR,
				levelStr: errorLevelStr,
				exitCode: -1,
			},
		},
		{
			name: "Warning",
			args: testLoggerLevelArgs{
				fn:  Warning,
				fnf: Warningf,
			},
			want: testLoggerLevelWant{
				level:    WARNING,
				levelStr: warningLevelStr,
				exitCode: -1,
			},
		},
		{
			name: "Info",
			args: testLoggerLevelArgs{
				fn:  Info,
				fnf: Infof,
			},
			want: testLoggerLevelWant{
				level:    INFO,
				levelStr: infoLevelStr,
				exitCode: -1,
			},
		},
		{
			name: "Debug",
			args: testLoggerLevelArgs{
				fn:  Debug,
				fnf: Debugf,
			},
			want: testLoggerLevelWant{
				level:    DEBUG,
				levelStr: debugLevelStr,
				exitCode: -1,
			},
		},
	}

	testLoggerLevels(t, std, testCases)
}
