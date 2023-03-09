package logger

// NewEncoderText creates a new text encoder.
func NewEncoderText(cfg EncoderTextConfig) *EncoderText {
	if cfg.Separator == "" {
		cfg.Separator = defaultTextSeparator
	}

	if cfg.DatetimeLayout == "" {
		cfg.DatetimeLayout = defaultDatetimeLayout
	}

	enc := new(EncoderText)
	enc.cfg = cfg

	return enc
}

// Copy returns a copy of the text encoder.
func (enc *EncoderText) Copy() Encoder {
	copyEnc := NewEncoderText(enc.cfg)
	copyEnc.EncoderBase = *enc.EncoderBase.Copy()

	return copyEnc
}

// Configure configures then encoder.
//
// - Encondes and sets the fields.
func (enc *EncoderText) Configure(cfg Config) {
	if len(cfg.Fields) == 0 {
		enc.SetFieldsEncoded("")

		return
	}

	buf := AcquireBuffer()

	for _, field := range cfg.Fields {
		buf.WriteString(field.Key) // nolint:errcheck
		buf.WriteString("=")       // nolint:errcheck
		buf.WriteInterface(field.Value)
		buf.WriteString(enc.cfg.Separator) // nolint:errcheck
	}

	enc.SetFieldsEncoded(buf.String())

	ReleaseBuffer(buf)
}

// Encode encodes the given entry to the buffer.
func (enc *EncoderText) Encode(buf *Buffer, e Entry) error {
	if e.Config.Datetime {
		buf.WriteDatetime(e.Time, enc.cfg.DatetimeLayout)
		buf.WriteString(enc.cfg.Separator) // nolint:errcheck
	}

	if e.Config.Timestamp {
		buf.WriteTimestamp(e.Time)
		buf.WriteString(enc.cfg.Separator) // nolint:errcheck
	}

	if levelStr := e.Level.String(); levelStr != "" {
		buf.WriteString(levelStr)          // nolint:errcheck
		buf.WriteString(enc.cfg.Separator) // nolint:errcheck
	}

	if e.Config.Shortfile || e.Config.Longfile {
		buf.WriteFileCaller(e.Caller, e.Config.Shortfile)
		buf.WriteString(enc.cfg.Separator) // nolint:errcheck
	}

	buf.WriteString(enc.FieldsEncoded()) // nolint:errcheck
	buf.WriteString(e.Message)           // nolint:errcheck
	buf.WriteNewLine()

	return nil
}
