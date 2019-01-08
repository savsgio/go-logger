package logger

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/valyala/bytebufferpool"
)

func TestLogger_isStd(t *testing.T) {
	tests := []struct {
		name string
		l    *Logger
		want bool
	}{
		{
			name: "Is standard",
			l:    New("std", "debug", &bytes.Buffer{}),
			want: true,
		},
		{
			name: "Is Not standard",
			l:    New("not", "debug", &bytes.Buffer{}),
			want: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.l.isStd(); got != test.want {
				t.Errorf("Logger.isStd() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestLogger_checkLevel(t *testing.T) {
	type args struct {
		level int
	}
	tests := []struct {
		name string
		l    *Logger
		args args
		want bool
	}{
		{
			name: "Fatal",
			l:    New("test", FATAL, &bytes.Buffer{}),
			args: args{level: fatalLevel},
			want: true,
		},
		{
			name: "Error",
			l:    New("test", ERROR, &bytes.Buffer{}),
			args: args{level: errorLevel},
			want: true,
		},
		{
			name: "Warning",
			l:    New("test", WARNING, &bytes.Buffer{}),
			args: args{level: warningLevel},
			want: true,
		},
		{
			name: "Info",
			l:    New("test", INFO, &bytes.Buffer{}),
			args: args{level: infoLevel},
			want: true,
		},
		{
			name: "Debug",
			l:    New("test", DEBUG, &bytes.Buffer{}),
			args: args{level: debugLevel},
			want: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.l.checkLevel(test.args.level); got != test.want {
				t.Errorf("Logger.checkLevel() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestLogger_writePrefix(t *testing.T) {
	type args struct {
		typeLevel string
	}

	tests := []struct {
		name string
		l    *Logger
		args args
		want string
	}{
		{
			name: "Standar",
			l:    New("std", DEBUG, &bytes.Buffer{}),
			args: args{typeLevel: "LevelString"},
			want: "- LevelString - ",
		},
		{
			name: "Not Standar",
			l:    New("not", DEBUG, &bytes.Buffer{}),
			args: args{typeLevel: "LevelString"},
			want: "- not - LevelString - ",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buff := bytebufferpool.Get()
			defer bytebufferpool.Put(buff)

			test.l.writePrefix(buff, test.args.typeLevel)
			got := buff.String()

			if got != test.want {
				t.Errorf("Logger.writePrefix() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestLogger_output(t *testing.T) {
	type args struct {
		prefix string
		msg    string
	}

	test := struct {
		args args
		want string
	}{
		args: args{prefix: "DEBUG", msg: "Test output"},
		want: " - DEBUG - Test output\n",
	}

	output := &bytes.Buffer{}
	l := New("test", DEBUG, output)

	l.output(test.args.prefix, test.args.msg)

	if got := output.String(); !strings.HasSuffix(got, test.want) {
		t.Errorf("Logger.output() = %v, want: %v", got, test.want)
	}

}

func TestLogger_outputf(t *testing.T) {
	type args struct {
		prefix string
		msg    string
		v      []interface{}
	}

	test := struct {
		args args
		want string
	}{
		args: args{prefix: "DEBUG", msg: "Test %s", v: []interface{}{"outputf"}},
		want: " - DEBUG - Test outputf\n",
	}

	output := &bytes.Buffer{}
	l := New("test", DEBUG, output)

	l.outputf(test.args.prefix, test.args.msg, test.args.v...)

	if got := output.String(); !strings.HasSuffix(got, test.want) {
		t.Errorf("Logger.outputf() = %v, want: %v", got, test.want)
	}
}

func TestLogger_SetLevel(t *testing.T) {
	type args struct {
		level    string
		intLevel int
	}
	tests := []struct {
		name     string
		args     args
		want     int
		getPanic bool
	}{
		{
			name:     "Fatal",
			args:     args{level: FATAL, intLevel: fatalLevel},
			want:     0,
			getPanic: false,
		},
		{
			name:     "Error",
			args:     args{level: ERROR, intLevel: errorLevel},
			want:     1,
			getPanic: false,
		},
		{
			name:     "Warning",
			args:     args{level: WARNING, intLevel: warningLevel},
			want:     2,
			getPanic: false,
		},
		{
			name:     "Info",
			args:     args{level: INFO, intLevel: infoLevel},
			want:     3,
			getPanic: false,
		},
		{
			name:     "Debug",
			args:     args{level: DEBUG, intLevel: debugLevel},
			want:     4,
			getPanic: false,
		},
		{
			name:     "Invalid",
			args:     args{level: "invalid"},
			want:     0,
			getPanic: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer func() {
				r := recover()

				if test.getPanic && r == nil {
					t.Errorf("Panic expected")
				} else if !test.getPanic && r != nil {
					t.Errorf("Unexpected panic")
				}
			}()

			l := New("test", FATAL, os.Stderr)
			l.SetLevel(test.args.level)

			if l.level != test.want {
				t.Errorf("Logger.SetLevel() = %v, want %v", l.level, test.want)
			}
		})
	}
}

func TestLogger_LevelEnabled(t *testing.T) {
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
	tests := []struct {
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
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := New("test", FATAL, os.Stderr)
			l.SetLevel(test.args.level)

			isEnabled := l.FatalEnabled()
			if isEnabled != test.want.fatalEnabled {
				t.Errorf("Logger.FatalEnabled() = '%v', want '%v'", isEnabled, test.want.fatalEnabled)
			}

			isEnabled = l.ErrorEnabled()
			if isEnabled != test.want.errorEnabled {
				t.Errorf("Logger.ErrorEnabled() = '%v', want '%v'", isEnabled, test.want.errorEnabled)
			}

			isEnabled = l.WarningEnabled()
			if isEnabled != test.want.warningEnabled {
				t.Errorf("Logger.WarningEnabled() = '%v', want '%v'", isEnabled, test.want.warningEnabled)
			}

			isEnabled = l.InfoEnabled()
			if isEnabled != test.want.infoEnabled {
				t.Errorf("Logger.InfoEnabled() = '%v', want '%v'", isEnabled, test.want.infoEnabled)
			}

			isEnabled = l.DebugEnabled()
			if isEnabled != test.want.debugEnabled {
				t.Errorf("Logger.DebugEnabled() = '%v', want '%v'", isEnabled, test.want.debugEnabled)
			}
		})
	}
}

func TestLogger_SetOutput(t *testing.T) {
	output := &bytes.Buffer{}

	l := New("std", DEBUG, os.Stdout)
	l.SetOutput(output)

	if l.out != output {
		t.Errorf("Logger.SetOutput() = %v, want %v", l.out, output)
	}
}

func TestLogger_Error(t *testing.T) {
	type args struct {
		msg []interface{}
	}

	tests := []struct {
		name    string
		logName string
		args    args
		want    string
	}{
		{
			name:    "Standar",
			logName: "std",
			args:    args{msg: []interface{}{"Msg with std"}},
			want:    "- ERROR - Msg with std\n",
		},
		{
			name:    "Not Standar",
			logName: "not",
			args:    args{msg: []interface{}{"Msg with not std"}},
			want:    "- not - ERROR - Msg with not std\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := &bytes.Buffer{}
			New(test.logName, ERROR, output).Error(test.args.msg...)

			if got := output.String(); !strings.HasSuffix(got, test.want) {
				t.Errorf("Logger.Error() = %v, want %v", got, test.want)
			}

		})
	}
}

func TestLogger_Errorf(t *testing.T) {
	type args struct {
		msg string
		v   []interface{}
	}

	tests := []struct {
		name    string
		logName string
		args    args
		want    string
	}{
		{
			name:    "Standar",
			logName: "std",
			args:    args{msg: "Msg with %s", v: []interface{}{"std"}},
			want:    "- ERROR - Msg with std\n",
		},
		{
			name:    "Not Standar",
			logName: "not",
			args:    args{msg: "Msg with %s", v: []interface{}{"not std"}},
			want:    "- not - ERROR - Msg with not std\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := &bytes.Buffer{}
			New(test.logName, INFO, output).Errorf(test.args.msg, test.args.v...)

			if got := output.String(); !strings.HasSuffix(got, test.want) {
				t.Errorf("Logger.Errorf() = %v, want %v", got, test.want)
			}

		})
	}
}

func TestLogger_Warning(t *testing.T) {
	type args struct {
		msg []interface{}
	}

	tests := []struct {
		name    string
		logName string
		args    args
		want    string
	}{
		{
			name:    "Standar",
			logName: "std",
			args:    args{msg: []interface{}{"Msg with std"}},
			want:    "- WARNING - Msg with std\n",
		},
		{
			name:    "Not Standar",
			logName: "not",
			args:    args{msg: []interface{}{"Msg with not std"}},
			want:    "- not - WARNING - Msg with not std\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := &bytes.Buffer{}
			New(test.logName, WARNING, output).Warning(test.args.msg...)

			if got := output.String(); !strings.HasSuffix(got, test.want) {
				t.Errorf("Logger.Warning() = %v, want %v", got, test.want)
			}

		})
	}
}

func TestLogger_Warningf(t *testing.T) {
	type args struct {
		msg string
		v   []interface{}
	}

	tests := []struct {
		name    string
		logName string
		args    args
		want    string
	}{
		{
			name:    "Standar",
			logName: "std",
			args:    args{msg: "Msg with %s", v: []interface{}{"std"}},
			want:    "- WARNING - Msg with std\n",
		},
		{
			name:    "Not Standar",
			logName: "not",
			args:    args{msg: "Msg with %s", v: []interface{}{"not std"}},
			want:    "- not - WARNING - Msg with not std\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := &bytes.Buffer{}
			New(test.logName, INFO, output).Warningf(test.args.msg, test.args.v...)

			if got := output.String(); !strings.HasSuffix(got, test.want) {
				t.Errorf("Logger.Warningf() = %v, want %v", got, test.want)
			}

		})
	}
}

func TestLogger_Info(t *testing.T) {
	type args struct {
		msg []interface{}
	}

	tests := []struct {
		name    string
		logName string
		args    args
		want    string
	}{
		{
			name:    "Standar",
			logName: "std",
			args:    args{msg: []interface{}{"Msg with std"}},
			want:    "- INFO - Msg with std\n",
		},
		{
			name:    "Not Standar",
			logName: "not",
			args:    args{msg: []interface{}{"Msg with not std"}},
			want:    "- not - INFO - Msg with not std\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := &bytes.Buffer{}
			New(test.logName, INFO, output).Info(test.args.msg...)

			if got := output.String(); !strings.HasSuffix(got, test.want) {
				t.Errorf("Logger.Info() = %v, want %v", got, test.want)
			}

		})
	}
}

func TestLogger_Infof(t *testing.T) {
	type args struct {
		msg string
		v   []interface{}
	}

	tests := []struct {
		name    string
		logName string
		args    args
		want    string
	}{
		{
			name:    "Standar",
			logName: "std",
			args:    args{msg: "Msg with %s", v: []interface{}{"std"}},
			want:    "- INFO - Msg with std\n",
		},
		{
			name:    "Not Standar",
			logName: "not",
			args:    args{msg: "Msg with %s", v: []interface{}{"not std"}},
			want:    "- not - INFO - Msg with not std\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := &bytes.Buffer{}
			New(test.logName, INFO, output).Infof(test.args.msg, test.args.v...)

			if got := output.String(); !strings.HasSuffix(got, test.want) {
				t.Errorf("Logger.Infof() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestLogger_Debug(t *testing.T) {
	type args struct {
		msg []interface{}
	}

	tests := []struct {
		name    string
		logName string
		args    args
		want    string
	}{
		{
			name:    "Standar",
			logName: "std",
			args:    args{msg: []interface{}{"Msg with std"}},
			want:    "- DEBUG - Msg with std\n",
		},
		{
			name:    "Not Standar",
			logName: "not",
			args:    args{msg: []interface{}{"Msg with not std"}},
			want:    "- not - DEBUG - Msg with not std\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := &bytes.Buffer{}
			New(test.logName, DEBUG, output).Debug(test.args.msg...)

			if got := output.String(); !strings.HasSuffix(got, test.want) {
				t.Errorf("Logger.Debug() = %v, want %v", got, test.want)
			}

		})
	}
}

func TestLogger_Debugf(t *testing.T) {
	type args struct {
		msg string
		v   []interface{}
	}

	tests := []struct {
		name    string
		logName string
		args    args
		want    string
	}{
		{
			name:    "Standar",
			logName: "std",
			args:    args{msg: "Msg with %s", v: []interface{}{"std"}},
			want:    "- DEBUG - Msg with std\n",
		},
		{
			name:    "Not Standar",
			logName: "not",
			args:    args{msg: "Msg with %s", v: []interface{}{"not std"}},
			want:    "- not - DEBUG - Msg with not std\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := &bytes.Buffer{}
			New(test.logName, DEBUG, output).Debugf(test.args.msg, test.args.v...)

			if got := output.String(); !strings.HasSuffix(got, test.want) {
				t.Errorf("Logger.Debugf() = %v, want %v", got, test.want)
			}

		})
	}
}

func TestLogger_Print(t *testing.T) {
	type args struct {
		msg []interface{}
	}

	tests := []struct {
		name    string
		logName string
		args    args
		want    string
	}{
		{
			name:    "Standar",
			logName: "std",
			args:    args{msg: []interface{}{"Msg with std"}},
			want:    "- Msg with std\n",
		},
		{
			name:    "Not Standar",
			logName: "not",
			args:    args{msg: []interface{}{"Msg with not std"}},
			want:    "- not - Msg with not std\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := &bytes.Buffer{}
			New(test.logName, DEBUG, output).Print(test.args.msg...)

			if got := output.String(); !strings.HasSuffix(got, test.want) {
				t.Errorf("Logger.Print() = %v, want %v", got, test.want)
			}

		})
	}
}

func TestLogger_Printf(t *testing.T) {
	type args struct {
		msg string
		v   []interface{}
	}

	tests := []struct {
		name    string
		logName string
		args    args
		want    string
	}{
		{
			name:    "Standar",
			logName: "std",
			args:    args{msg: "Msg with %s", v: []interface{}{"std"}},
			want:    "Msg with std\n",
		},
		{
			name:    "Not Standar",
			logName: "not",
			args:    args{msg: "Msg with %s", v: []interface{}{"not std"}},
			want:    "Msg with not std\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := &bytes.Buffer{}
			New(test.logName, DEBUG, output).Printf(test.args.msg, test.args.v...)

			if got := output.String(); !strings.HasSuffix(got, test.want) {
				t.Errorf("Logger.Printf() = %v, want %v", got, test.want)
			}

		})
	}
}

// Benchmarks
func Benchmark_Printf(b *testing.B) {
	output := &bytes.Buffer{}
	log := New("test", DEBUG, output)

	b.ResetTimer()
	for i := 0; i <= b.N; i++ {
		log.Printf("Test %s", "params")
	}
}

func Benchmark_Errorf(b *testing.B) {
	output := &bytes.Buffer{}
	log := New("test", DEBUG, output)

	b.ResetTimer()
	for i := 0; i <= b.N; i++ {
		log.Errorf("Test %s", "params")
	}
}

func Benchmark_Warningf(b *testing.B) {
	output := &bytes.Buffer{}
	log := New("test", DEBUG, output)

	b.ResetTimer()
	for i := 0; i <= b.N; i++ {
		log.Warningf("Test %s", "params")
	}
}

func Benchmark_Infof(b *testing.B) {
	output := &bytes.Buffer{}
	log := New("test", DEBUG, output)

	b.ResetTimer()
	for i := 0; i <= b.N; i++ {
		log.Infof("Test %s", "params")
	}
}

func Benchmark_Debugf(b *testing.B) {
	output := &bytes.Buffer{}
	log := New("test", DEBUG, output)

	b.ResetTimer()
	for i := 0; i <= b.N; i++ {
		log.Debugf("Test %s", "params")
	}
}
