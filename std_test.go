package logger

import (
	"sync"
	"testing"
)

var stdLock = sync.Mutex{}

func acquireStd() {
	stdLock.Lock()
}

func releaseStd() {
	std = newStd()
	stdLock.Unlock()
}

func TestLogger_std_WithFields(t *testing.T) {
	acquireStd()

	defer releaseStd()

	testLoggerWithFields(t, std, WithFields)
}

func TestLogger_std_SetFields(t *testing.T) {
	acquireStd()

	defer releaseStd()

	testLoggerSetFields(t, std, SetFields)
}

func TestLogger_std_SetFlags(t *testing.T) {
	acquireStd()

	defer releaseStd()

	testLoggerSetFlags(t, std, SetFlags)
}

func TestLogger_std_SetLevel(t *testing.T) {
	acquireStd()

	defer releaseStd()

	testLoggerSetLevel(t, std, SetLevel)
}

func TestLogger_std_SetOutput(t *testing.T) {
	acquireStd()

	defer releaseStd()

	testLoggerSetOutput(t, std, SetOutput)
}

func TestLogger_std_SetEncoder(t *testing.T) {
	acquireStd()

	defer releaseStd()

	testLoggerSetEncoder(t, std, SetEncoder)
}

func TestLogger_std_IsLevelEnabled(t *testing.T) {
	acquireStd()

	defer releaseStd()

	testLoggerIsLevelEnabled(t, std, IsLevelEnabled)
}

func TestLogger_std_AddHook(t *testing.T) {
	acquireStd()

	defer releaseStd()

	testLoggerAddHook(t, std, AddHook)
}

func TestLogger_std_Levels(t *testing.T) { // nolint:funlen
	acquireStd()

	defer releaseStd()

	std = newTestLogger()
	std.setCalldepth(calldepthStd)

	testCases := []testLoggerLevelCase{
		{
			name: "Print",
			args: testLoggerLevelArgs{
				fn:  Print,
				fnf: Printf,
			},
			want: testLoggerLevelWant{
				level:    PRINT,
				exitCode: -1,
			},
		},
		{
			name: "Panic",
			args: testLoggerLevelArgs{
				fn:  Panic,
				fnf: Panicf,
			},
			want: testLoggerLevelWant{
				level:    PANIC,
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
				exitCode: -1,
			},
		},
	}

	testLoggerLevels(t, std, testCases)
}
