package logger

func newEncoderBase() *EncoderBase {
	return new(EncoderBase)
}

// Copy returns a copy of the encoder base.
func (enc *EncoderBase) Copy() *EncoderBase {
	copyEnc := newEncoderBase()
	copyEnc.fieldsEncoded = enc.fieldsEncoded

	return copyEnc
}

// FieldsEnconded returns the encoded fields.
func (enc *EncoderBase) FieldsEnconded() string {
	return enc.fieldsEncoded
}

// SetFieldsEnconded sets the fields enconded.
func (enc *EncoderBase) SetFieldsEnconded(fieldsEncoded string) {
	enc.fieldsEncoded = fieldsEncoded
}
