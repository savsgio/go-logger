package logger

import (
	"io"
	"log"
	"sync"
)

type Logger struct {
	mu       sync.Mutex // ensures atomic writes; protects the following fields
	name     string
	level    int
	out      io.Writer
	instance *log.Logger
	flag     int

	fatalEnabled   bool
	errorEnabled   bool
	warningEnabled bool
	infoEnabled    bool
	debugEnabled   bool
}
