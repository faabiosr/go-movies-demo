package errors

import (
	"net/http"
	"testing"
)

func TestHTTP(t *testing.T) {
	tests := []struct {
		name string
		err  interface{}
		want string
	}{
		{
			name: "no error",
		},
		{
			name: "error",
			err:  New("validation failed"),
			want: "validation failed",
		},
		{
			name: "string",
			err:  "validation failed",
			want: "validation failed",
		},
		{
			name: "panic",
			err:  0,
			want: "err must be error or string, got int",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if r != nil && r != tt.want {
					t.Errorf("want %v, got %v", tt.want, r)
				}
			}()

			err := HTTP(http.StatusBadRequest, tt.err)
			if err != nil && tt.want != err.Error() {
				t.Errorf("want %v, got %v", tt.want, err.Error())
			}
		})
	}

	t.Run("http error", func(t *testing.T) {
		err := New("validation failed")
		got := HTTP(http.StatusBadRequest, err)

		if code := got.(*httpErr).HTTPStatus(); code != http.StatusBadRequest {
			t.Errorf("status code: want %v, got %v", http.StatusBadRequest, code)
		}

		if !Is(got, err) {
			t.Errorf("error is not the same: want %v, got %v", err, got)
		}
	})
}
