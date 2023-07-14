package movies

import (
	"testing"
)

func TestDataSourceStore(t *testing.T) {
	ds := NewDatasource(db(t))

	m := Movie{
		Title:    "Star Wars: Episode I - The Phantom Menace",
		Released: newDate(t, "1999-08-20"),
	}

	_, err := ds.Store(m)
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestDataSourceFetch(t *testing.T) {
	t.Run("fails", func(t *testing.T) {
		ds := NewDatasource(db(t))

		_, err := ds.Fetch("123456")
		if err == nil {
			t.Error("expected an error, got nil")
		}
	})

	t.Run("success", func(t *testing.T) {
		ds := NewDatasource(db(t))

		m := Movie{
			Title:    "Star Wars: Episode II - Attack of the clones",
			Released: newDate(t, "2002-05-17"),
		}

		m, _ = ds.Store(m)

		got, err := ds.Fetch(m.ID)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}

		if got.Title != m.Title {
			t.Errorf("expected %s, got %s", m.Title, got.Title)
		}
	})
}

func TestDataSourceFetchAll(t *testing.T) {
	ds := NewDatasource(db(t))

	m := Movie{
		Title:    "Star Wars: Episode II - Attack of the clones",
		Released: newDate(t, "2002-05-17"),
	}

	_, _ = ds.Store(m)

	movies, err := ds.FetchAll()
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	if total := len(movies); total == 0 {
		t.Errorf("expected at least 1 movie, got %d", total)
	}
}

func TestDatasourceRemove(t *testing.T) {
	t.Run("empty id", func(t *testing.T) {
		ds := NewDatasource(db(t))

		if err := ds.Remove(""); err == nil {
			t.Error("expected an error, got nil")
		}
	})

	t.Run("not found", func(t *testing.T) {
		ds := NewDatasource(db(t))

		if err := ds.Remove("1"); err == nil {
			t.Error("expected an error, got nil")
		}
	})

	t.Run("success", func(t *testing.T) {
		ds := NewDatasource(db(t))

		m := Movie{
			Title:    "Star Wars: Episode II - Attack of the clones",
			Released: newDate(t, "2002-05-17"),
		}

		m, _ = ds.Store(m)

		if err := ds.Remove(m.ID); err != nil {
			t.Errorf("expected nil, got %v", err)
		}

		movies, _ := ds.FetchAll()
		if total := len(movies); total != 0 {
			t.Errorf("expected not movies, got %d", total)
		}
	})
}
