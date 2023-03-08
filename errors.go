package logger

import "errors"

var (
	// ErrInvalidLevel is the invalid level error.
	ErrInvalidLevel = errors.New("invalid level")

	// ErrEmptyHookLevels is the empty hook levels error.
	ErrEmptyHookLevels = errors.New("empty hook levels")
)
