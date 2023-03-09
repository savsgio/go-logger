package logger

// Copy returns a copy of the config.
func (e Config) Copy() Config {
	e2 := e

	if e.Fields != nil {
		e2.Fields = make([]Field, len(e.Fields))
		copy(e2.Fields, e.Fields)
	}

	return e2
}
