package errors

type publicErr struct {
	err error
	msg string
}

func (e *publicErr) Error() string {
	return e.err.Error()
}

func (e *publicErr) Public() string {
	return e.msg
}

func (e *publicErr) Unwrap() error {
	return e.err
}

// Public wraps an error with a public error message. The public error message
// must be explicitly extrated, as so:
//
//	var pe interface {
//	    Public() string
//	}
//
//	if errors.As(err, pe); ok {
//	    fmt.Println(pe.Public())
//	}
//
// Otherwise, it behaves as standard `err` simply passed through.
func Public(err error, msg string) error {
	if err == nil {
		return nil
	}

	return &publicErr{
		err: err,
		msg: msg,
	}
}
