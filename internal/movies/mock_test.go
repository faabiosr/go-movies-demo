package movies

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	bolt "go.etcd.io/bbolt"
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

func db(t *testing.T) *bolt.DB {
	t.Helper()

	dir, err := os.MkdirTemp("", strings.ReplaceAll(t.Name(), "/", "_"))
	if err != nil {
		t.Fatal(err)
	}

	dbPath := filepath.Join(dir, "test.db")
	db, err := bolt.Open(dbPath, 0o600, nil)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		_ = db.Close()
		_ = os.RemoveAll(dir)
	})

	return db
}
