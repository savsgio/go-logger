package logger

import (
	"errors"
	"testing"
)

func Test_ParseLevel(t *testing.T) { // nolint:funlen
	type args struct {
		levelStr string
	}

	type want struct {
		level Level
		err   error
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Print",
			args: args{
				levelStr: "",
			},
			want: want{
				level: PRINT,
				err:   nil,
			},
		},
		{
			name: "Trace",
			args: args{
				levelStr: "trace",
			},
			want: want{
				level: TRACE,
				err:   nil,
			},
		},
		{
			name: "Panic",
			args: args{
				levelStr: "paNiC",
			},
			want: want{
				level: PANIC,
				err:   nil,
			},
		},
		{
			name: "Fatal",
			args: args{
				levelStr: "FatAL",
			},
			want: want{
				level: FATAL,
				err:   nil,
			},
		},
		{
			name: "Error",
			args: args{
				levelStr: "error",
			},
			want: want{
				level: ERROR,
				err:   nil,
			},
		},
		{
			name: "Warning",
			args: args{
				levelStr: "WarNIng",
			},
			want: want{
				level: WARNING,
				err:   nil,
			},
		},
		{
			name: "Info",
			args: args{
				levelStr: "INFO",
			},
			want: want{
				level: INFO,
				err:   nil,
			},
		},
		{
			name: "Debug",
			args: args{
				levelStr: "dEBUg",
			},
			want: want{
				level: DEBUG,
				err:   nil,
			},
		},
		{
			name: "Invalid",
			args: args{
				levelStr: "inValId",
			},
			want: want{
				level: invalid,
				err:   ErrInvalidLevel,
			},
		},
	}

	for i := range tests {
		test := tests[i]

		t.Run(test.name, func(t *testing.T) {
			t.Helper()

			level, err := ParseLevel(test.args.levelStr)

			if level != test.want.level {
				t.Errorf("level == %d, want %d", level, test.want.level)
			}

			if !errors.Is(err, test.want.err) {
				t.Errorf("error == %d, want %d", err, test.want.err)
			}
		})
	}
}

func TestLevel_String(t *testing.T) { // nolint:funlen
	type args struct {
		level Level
	}

	type want struct {
		result string
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Print",
			args: args{
				level: PRINT,
			},
			want: want{
				result: printLevelStr,
			},
		},
		{
			name: "Trace",
			args: args{
				level: TRACE,
			},
			want: want{
				result: traceLevelStr,
			},
		},
		{
			name: "Panic",
			args: args{
				level: PANIC,
			},
			want: want{
				result: panicLevelStr,
			},
		},
		{
			name: "Fatal",
			args: args{
				level: FATAL,
			},
			want: want{
				result: fatalLevelStr,
			},
		},
		{
			name: "Error",
			args: args{
				level: ERROR,
			},
			want: want{
				result: errorLevelStr,
			},
		},
		{
			name: "Warning",
			args: args{
				level: WARNING,
			},
			want: want{
				result: warningLevelStr,
			},
		},
		{
			name: "Info",
			args: args{
				level: INFO,
			},
			want: want{
				result: infoLevelStr,
			},
		},
		{
			name: "Debug",
			args: args{
				level: DEBUG,
			},
			want: want{
				result: debugLevelStr,
			},
		},
		{
			name: "Invalid",
			args: args{
				level: invalid,
			},
			want: want{
				result: ErrInvalidLevel.Error(),
			},
		},
	}

	for i := range tests {
		test := tests[i]

		t.Run(test.name, func(t *testing.T) {
			t.Helper()

			if result := test.args.level.String(); result != test.want.result {
				t.Errorf("level == %s, want %s", result, test.want.result)
			}
		})
	}
}
