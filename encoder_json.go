package logger

// NewEncoderJSON creates a new json encoder.
func NewEncoderJSON() *EncoderJSON {
	return new(EncoderJSON)
}

// Copy returns a copy of the json encoder.
func (enc *EncoderJSON) Copy() Encoder {
	copyEnc := NewEncoderJSON()
	copyEnc.EncoderBase = *enc.EncoderBase.Copy()

	return copyEnc
}

// SetFields encodes and sets the given fields.
func (enc *EncoderJSON) SetFields(fields []Field) {
	buf := AcquireBuffer()

	for _, field := range fields {
		buf.WriteString("\"")      // nolint:errcheck
		buf.WriteString(field.Key) // nolint:errcheck
		buf.WriteString("\":\"")   // nolint:errcheck

		n := buf.Len()
		buf.WriteInterface(field.Value)
		buf.Escape(n)

		buf.WriteString("\",") // nolint:errcheck
	}

	enc.SetFieldsEncoded(buf.String())

	ReleaseBuffer(buf)
}

// Encode encodes the given entry to the buffer.
func (enc *EncoderJSON) Encode(buf *Buffer, e Entry) error {
	buf.WriteByte('{') // nolint:errcheck

	if e.Config.Datetime {
		buf.WriteString("\"datetime\":\"") // nolint:errcheck
		buf.WriteDatetime(e.Time)
		buf.WriteString("\",") // nolint:errcheck
	}

	if e.Config.Timestamp {
		buf.WriteString("\"timestamp\":\"") // nolint:errcheck
		buf.WriteTimestamp(e.Time)
		buf.WriteString("\",") // nolint:errcheck
	}

	if levelStr := e.Level.String(); levelStr != "" {
		buf.WriteString("\"level\":\"") // nolint:errcheck
		buf.WriteString(levelStr)       // nolint:errcheck
		buf.WriteString("\",")          // nolint:errcheck
	}

	if e.Config.Shortfile || e.Config.Longfile {
		buf.WriteString("\"file\":\"") // nolint:errcheck
		buf.WriteFileCaller(e.Caller, e.Config.Shortfile)
		buf.WriteString("\",") // nolint:errcheck
	}

	buf.WriteString(enc.FieldsEncoded()) // nolint:errcheck
	buf.WriteString("\"message\":\"")    // nolint:errcheck

	n := buf.Len()
	buf.WriteString(e.Message) // nolint:errcheck
	buf.Escape(n)

	buf.WriteString("\"}") // nolint:errcheck
	buf.WriteNewLine()

	return nil
}
