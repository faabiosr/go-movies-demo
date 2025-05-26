package errors

import (
	"fmt"
)

type httpErr struct {
	code int
	err  error
}

func (e *httpErr) Error() string {
	return e.err.Error()
}

func (e *httpErr) Unwrap() error {
	return e.err
}

func (e *httpErr) HTTPStatus() int {
	return e.code
}

// HTTP wraps err with the specified HTTP status code. err could be either an
// error or a string, othewise panic. If err is a nil error, the return value
// will also be nil.
func HTTP(status int, err any) error {
	if err == nil {
		return nil
	}

	var e error
	switch t := err.(type) {
	case error:
		e = t
	case string:
		e = New(t)
	default:
		panic(fmt.Sprintf("err must be error or string, got %T", err))
	}

	return &httpErr{
		code: status,
		err:  e,
	}
}
