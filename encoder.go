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

// FieldsEncoded returns the encoded fields.
func (enc *EncoderBase) FieldsEncoded() string {
	return enc.fieldsEncoded
}

// SetFieldsEncoded sets the fields enconded.
func (enc *EncoderBase) SetFieldsEncoded(fieldsEncoded string) {
	enc.fieldsEncoded = fieldsEncoded
}
