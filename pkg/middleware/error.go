package middleware

import (
	"fmt"
	"net/http"

	"github.com/invopop/validation"
	"github.com/labstack/echo/v4"

	"github.com/faabiosr/go-movies-demo/pkg/errors"
)

// ErrorHandler handles the default echo error to dispatch in json format.
func ErrorHandler(err error, ec echo.Context) {
	if err == nil {
		return
	}

	status, msg := inspectErr(err)

	if ec.Response().Committed {
		return
	}

	if ec.Request().Method == http.MethodHead {
		_ = ec.NoContent(status)
		return
	}

	_ = ec.JSON(status, echo.Map{"message": msg})
}

func inspectErr(err error) (int, string) {
	return status(err), message(err)
}

func status(err error) int {
	var he interface {
		HTTPStatus() int
	}

	if errors.As(err, &he) {
		return he.HTTPStatus()
	}

	var ve validation.Errors

	if errors.As(err, &ve) {
		return http.StatusBadRequest
	}

	if ee := new(echo.HTTPError); errors.As(err, &ee) {
		return ee.Code
	}

	return http.StatusInternalServerError
}

func message(err error) string {
	var pe interface {
		Public() string
	}

	if errors.As(err, &pe) {
		return pe.Public()
	}

	if ee := new(echo.HTTPError); errors.As(err, &ee) {
		return fmt.Sprintf("%v", ee.Message)
	}

	return err.Error()
}
