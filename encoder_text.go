package logger

const sepText = " - "

// NewEncoderText creates a new text encoder.
func NewEncoderText() *EncoderText {
	return new(EncoderText)
}

// Copy returns a copy of the text encoder.
func (enc *EncoderText) Copy() Encoder {
	copyEnc := NewEncoderText()
	copyEnc.EncoderBase = *enc.EncoderBase.Copy()

	return copyEnc
}

// SetFields encodes and sets the given fields.
func (enc *EncoderText) SetFields(fields []Field) {
	if len(fields) == 0 {
		enc.SetFieldsEncoded("")

		return
	}

	buf := AcquireBuffer()
	buf.WriteString("{") // nolint:errcheck

	for i, field := range fields {
		if i > 0 {
			buf.WriteString(",") // nolint:errcheck
		}

		buf.WriteString("\"")      // nolint:errcheck
		buf.WriteString(field.Key) // nolint:errcheck
		buf.WriteString("\":\"")   // nolint:errcheck
		buf.WriteInterface(field.Value)
		buf.WriteString("\"") // nolint:errcheck
	}

	buf.WriteString("}")     // nolint:errcheck
	buf.WriteString(sepText) // nolint:errcheck

	enc.SetFieldsEncoded(buf.String())

	ReleaseBuffer(buf)
}

// Encode encodes the given entry to the buffer.
func (enc *EncoderText) Encode(buf *Buffer, e Entry) error {
	if e.Config.Datetime {
		buf.WriteDatetime(e.Time)
		buf.WriteString(sepText) // nolint:errcheck
	}

	if e.Config.Timestamp {
		buf.WriteTimestamp(e.Time)
		buf.WriteString(sepText) // nolint:errcheck
	}

	if levelStr := e.Level.String(); levelStr != "" {
		buf.WriteString(levelStr) // nolint:errcheck
		buf.WriteString(sepText)  // nolint:errcheck
	}

	if e.Config.Shortfile || e.Config.Longfile {
		buf.WriteFileCaller(e.Caller, e.Config.Shortfile)
		buf.WriteString(sepText) // nolint:errcheck
	}

	buf.WriteString(enc.FieldsEncoded()) // nolint:errcheck
	buf.WriteString(e.Message)           // nolint:errcheck
	buf.WriteNewLine()

	return nil
}
