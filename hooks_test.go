package logger

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"testing"
	"time"
)

type testHook struct {
	levels   []Level
	fireFunc func(Entry) error
}

func (h *testHook) Levels() []Level {
	return h.levels
}

func (h *testHook) Fire(e Entry) error {
	return h.fireFunc(e)
}

func levelInclude(values []Level, level Level) bool {
	for i := range values {
		if values[i] == level {
			return true
		}
	}

	return false
}

func levelsToInts(levels []Level) (result []int) {
	for _, l := range levels {
		result = append(result, int(l))
	}

	sort.Ints(result)

	return result
}

func TestLevelHooks_copy(t *testing.T) {
	lh := newLevelHooks()

	if err := lh.add(&testHook{levels: []Level{DEBUG, ERROR}}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lh2 := lh.copy()

	if reflect.ValueOf(lh).Pointer() == reflect.ValueOf(lh2).Pointer() {
		t.Errorf("same pointers")
	}

	if reflect.ValueOf(lh.store).Pointer() == reflect.ValueOf(lh2.store).Pointer() {
		t.Errorf("store has the same pointers")
	}

	if !reflect.DeepEqual(lh.store, lh2.store) {
		t.Errorf("store values are not equals")
	}

	if lh.errOutput != lh2.errOutput {
		t.Errorf("error outputs are not equals")
	}
}

func TestLevelHooks_add(t *testing.T) { // nolint:funlen
	type args struct {
		hook *testHook
	}

	type want struct {
		err error
	}

	tests := []struct {
		args args
		want want
	}{
		{
			args: args{
				hook: &testHook{
					levels: []Level{INFO, DEBUG},
				},
			},
			want: want{
				err: nil,
			},
		},
		{
			args: args{
				hook: &testHook{
					levels: []Level{},
				},
			},
			want: want{
				err: ErrEmptyHookLevels,
			},
		},
	}

	for i := range tests {
		test := tests[i]

		t.Run("", func(t *testing.T) {
			lh := newLevelHooks()

			if err := lh.add(test.args.hook); !errors.Is(err, test.want.err) {
				t.Errorf("error == %v, want %v", err, test.want.err)
			}

			resultLevels := make([]Level, 0)

			for level, hooks := range lh.store {
				resultLevels = append(resultLevels, level)

				for _, h := range hooks {
					if reflect.ValueOf(h).Pointer() != reflect.ValueOf(test.args.hook).Pointer() {
						t.Errorf("hook == %p, want %p", h, test.args.hook)
					}
				}
			}

			resultLevelsInts := levelsToInts(resultLevels)
			hookLevels := levelsToInts(test.args.hook.levels)

			if !reflect.DeepEqual(resultLevelsInts, hookLevels) {
				t.Errorf("hook levels == %v, want %v", resultLevels, test.args.hook.levels)
			}
		})
	}
}

func TestLevelHooks_fire(t *testing.T) { // nolint:funlen
	type args struct {
		hook *testHook
	}

	type want struct {
		err error
	}

	hookFired := false
	hookEntry := Entry{}

	fireFunc := func(e Entry) error { // nolint:unparam
		hookFired = true
		hookEntry = e

		return nil
	}

	someErr := errors.New("some error")

	tests := []struct {
		args args
		want want
	}{
		{
			args: args{
				hook: &testHook{
					levels:   []Level{INFO, DEBUG},
					fireFunc: fireFunc,
				},
			},
			want: want{
				err: nil,
			},
		},
		{
			args: args{
				hook: &testHook{
					levels: []Level{INFO, DEBUG},
					fireFunc: func(e Entry) error {
						_ = fireFunc(e)

						return someErr
					},
				},
			},
			want: want{
				err: someErr,
			},
		},
	}

	for i := range tests {
		test := tests[i]

		errOutput := bytes.NewBuffer(nil)

		lh := newLevelHooks()
		lh.errOutput = errOutput

		if err := lh.add(test.args.hook); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		for _, level := range levels {
			t.Run(level.String(), func(t *testing.T) {
				e := Entry{
					Level:   level,
					Time:    time.Now(),
					Message: fmt.Sprintf("Hello %s", level),
				}

				expectedFired := levelInclude(test.args.hook.levels, level)

				lh.fire(e)

				if hookFired != expectedFired {
					t.Errorf("fired == %t, want %t", hookFired, expectedFired)
				}

				if expectedFired && !reflect.DeepEqual(hookEntry, e) {
					t.Errorf("entry == %v, want %v", hookEntry, e)
				}

				expectedErrorMsg := ""
				if expectedFired && test.want.err != nil {
					expectedErrorMsg = fmt.Sprintf("failed to fire hook[%s][%d]: %+v\n", e.Level, 0, test.want.err)
				}

				if errorMsg := errOutput.String(); errorMsg != expectedErrorMsg {
					t.Errorf("error message == %s, want %s", errorMsg, expectedErrorMsg)
				}

				hookFired = false
				hookEntry = Entry{}

				errOutput.Reset()
			})
		}
	}
}
