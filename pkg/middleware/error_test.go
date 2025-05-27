package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/invopop/validation"
	"github.com/labstack/echo/v4"

	"github.com/faabiosr/go-movies-demo/pkg/errors"
)

func TestErrorHandler(t *testing.T) {
	tests := []struct {
		name    string
		method  string
		err     error
		status  int
		content string
	}{
		{
			name:   "no error",
			err:    nil,
			status: http.StatusOK,
		},
		{
			name:    "standard error",
			err:     errors.New("standard"),
			status:  http.StatusInternalServerError,
			content: `{"message":"standard"}`,
		},
		{
			name:   "no content",
			method: http.MethodHead,
			err:    errors.New("foo"),
			status: http.StatusInternalServerError,
		},
		{
			name:    "echo error",
			err:     &echo.HTTPError{Code: http.StatusBadRequest, Message: "echo error"},
			status:  http.StatusBadRequest,
			content: `{"message":"echo error"}`,
		},
		{
			name:    "http error",
			err:     errors.HTTP(http.StatusNotFound, "event was not found"),
			status:  http.StatusNotFound,
			content: `{"message":"event was not found"}`,
		},

		{
			name:    "public error",
			err:     errors.Public(errors.New("failed"), "unable to store"),
			status:  http.StatusInternalServerError,
			content: `{"message":"unable to store"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.method == "" {
				tt.method = http.MethodGet
			}

			req := httptest.NewRequest(tt.method, "/", nil)
			rec := httptest.NewRecorder()
			ec := echo.New().NewContext(req, rec)

			ErrorHandler(tt.err, ec)

			res := rec.Result()

			if code := res.StatusCode; code != tt.status {
				t.Errorf("expected status code '%d', got %d", tt.status, code)
			}

			body, err := io.ReadAll(res.Body)
			if err != nil {
				t.Error(err)
			}

			content := strings.TrimSpace(string(body))

			if strings.Compare(content, tt.content) != 0 {
				t.Errorf("expected message '%s', got '%s'", tt.content, content)
			}
		})
	}

	t.Run("committed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		ec := echo.New().NewContext(req, rec)
		ec.Response().WriteHeader(http.StatusBadRequest)

		ErrorHandler(errors.New("failed"), ec)

		res := rec.Result()

		if code := res.StatusCode; code != http.StatusBadRequest {
			t.Errorf("expected status code '%d', got %d", http.StatusBadRequest, code)
		}

		if res.ContentLength >= 0 {
			t.Error("expected empty content")
		}
	})

	t.Run("validation error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		ec := echo.New().NewContext(req, rec)

		c := struct {
			Name string
		}{}
		ErrorHandler(validation.Errors{
			"name": validation.Validate(c.Name, validation.Required),
		}, ec)

		res := rec.Result()

		if code := res.StatusCode; code != http.StatusBadRequest {
			t.Errorf("expected status code '%d', got %d", http.StatusBadRequest, code)
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			t.Error(err)
		}

		content := strings.TrimSpace(string(body))
		expected := `{"message":"name: cannot be blank."}`

		if strings.Compare(content, expected) != 0 {
			t.Errorf("expected message '%s', got '%s'", expected, content)
		}
	})
}
