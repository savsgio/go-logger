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
