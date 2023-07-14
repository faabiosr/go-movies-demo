package movies

import (
	"testing"
	"time"
)

func newDate(t *testing.T, v string) Date {
	t.Helper()

	return Date{
		Time: func() time.Time {
			tn, _ := time.Parse("2006-01-02", v)
			return tn
		}(),
	}
}
