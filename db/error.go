package db

type ValidateError struct {
	msg string
}

func (err *ValidateError) Error() string {
	return err.msg
}

func newValidateError(s string) *ValidateError {
	err := new(ValidateError)
	err.msg = s
	return err
}
