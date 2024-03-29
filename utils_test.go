package logger

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func Test_getFileCaller(t *testing.T) { // nolint:funlen
	type args struct {
		calldepth int
	}

	type want struct {
		frame runtime.Frame
	}

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			args: args{
				calldepth: 2,
			},
			want: want{
				frame: runtime.Frame{
					File: filepath.Join(cwd, "utils_test.go"),
					Line: 60,
				},
			},
		},
		{
			name: "invalid calldepth",
			args: args{
				calldepth: 1000,
			},
			want: want{
				frame: runtime.Frame{
					File: "???",
					Line: 0,
				},
			},
		},
	}

	for i := range tests {
		test := tests[i]

		t.Run(test.name, func(t *testing.T) {
			t.Helper()

			frame := getFileCaller(test.args.calldepth)

			frameFile := filepath.ToSlash(frame.File)
			wantFile := filepath.ToSlash(test.want.frame.File)

			if frameFile != wantFile {
				t.Errorf("file == %s, want %s", frameFile, wantFile)
			}

			if frame.Line != test.want.frame.Line {
				t.Errorf("line == %d, want %d", frame.Line, test.want.frame.Line)
			}
		})
	}
}
