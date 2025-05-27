package movies

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestRouterCreate(t *testing.T) {
	tests := []struct {
		name string
		body string
		err  string
	}{
		{
			name: "invalid json",
			body: "{",
			err:  "code=400, message=unexpected EOF, internal=unexpected EOF",
		},
		{
			name: "validation failed",
			body: "{}",
			err:  "released: cannot be blank; title: cannot be blank.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/movies", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e := echo.New()
			ec := e.NewContext(req, rec)

			r := &router{}

			err := r.create(ec)
			if err.Error() != tt.err {
				t.Errorf("unexpected error: %s (expected %s)", err.Error(), tt.err)
			}
		})
	}
}
