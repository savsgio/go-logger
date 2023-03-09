package logger

import (
	"fmt"
	"os"
)

func newLevelHooks() *levelHooks {
	return &levelHooks{
		store:     make(map[Level][]Hook),
		errOutput: os.Stderr,
	}
}

func (lh levelHooks) copy() *levelHooks {
	lh2 := newLevelHooks()
	lh2.errOutput = lh.errOutput

	for level, hooks := range lh.store {
		lh2.store[level] = append(lh2.store[level], hooks...)
	}

	return lh2
}

func (lh levelHooks) add(h Hook) error {
	levels := h.Levels()

	if len(levels) == 0 {
		return ErrEmptyHookLevels
	}

	for _, level := range levels {
		lh.store[level] = append(lh.store[level], h)
	}

	return nil
}

func (lh levelHooks) fire(e Entry) {
	hooks := lh.store[e.Level]

	for i := range hooks {
		if err := hooks[i].Fire(e); err != nil {
			fmt.Fprintf(lh.errOutput, "failed to fire hook[%s][%d]: %+v\n", e.Level, i, err)
		}
	}
}
