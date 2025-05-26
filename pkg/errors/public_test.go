package errors

import (
	"testing"
)

func TestPublic(t *testing.T) {
	err := New("failed to connect")
	perr := Public(err, "unable to access")

	t.Run("ignore nil", func(t *testing.T) {
		if err := Public(nil, "error"); err != nil {
			t.Errorf("expected nil, got %s", err)
		}
	})

	t.Run("passed through", func(t *testing.T) {
		if perr.Error() != err.Error() {
			t.Errorf("expected %s, got %s", err, perr)
		}
	})

	t.Run("wrapped", func(t *testing.T) {
		if got := Unwrap(perr); err != got {
			t.Errorf("expected %T, got %T", err, got)
		}
	})

	t.Run("public", func(t *testing.T) {
		var pe interface {
			Public() string
		}

		ok := As(perr, &pe)
		if !ok {
			t.Error("invalid error type")
		}

		if expect := "unable to access"; pe.Public() != expect {
			t.Errorf("expected %s, got %s", expect, pe.Public())
		}
	})
}
