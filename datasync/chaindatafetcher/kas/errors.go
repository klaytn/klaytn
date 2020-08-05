package kas

import "errors"

var (
	noOpcodeError    = errors.New("trace result must include type field")
	noFromFieldError = errors.New("from field is missing")
	noToFieldError   = errors.New("to field is missing")
)
