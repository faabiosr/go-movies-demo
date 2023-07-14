package movies

import (
	"encoding/json"
	"testing"
)

func TestDateMarshalJSON(t *testing.T) {
	date := newDate(t, "2023-07-13")

	s, err := json.Marshal(date)
	if err != nil {
		t.Errorf("marshal failed: expected nil, got %v", err)
	}

	if string(s) != `"2023-07-13"` {
		t.Error("date was not correctly formatted")
	}
}

func TestDateUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		json string
		err  string
	}{
		{
			name: "empty",
			json: `{"date": ""}`,
			err:  "",
		},
		{
			name: "invalid format",
			json: `{"date": "20-02-2002"}`,
			err:  `parsing time "20-02-2002" as "2006-01-02": cannot parse "20-02-2002" as "2006"`,
		},
		{
			name: "valid",
			json: `{"date": "2022-02-20"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var res struct {
				Date Date `json:"date"`
			}

			err := json.Unmarshal([]byte(tt.json), &res)
			if err == nil {
				return
			}

			if err.Error() != tt.err {
				t.Errorf("unexpected error: %s, got %s", tt.err, err.Error())
			}
		})
	}
}

func TestMovieValidation(t *testing.T) {
	m := Movie{}

	if err := m.Validate(); err == nil {
		t.Error("expected an error, got nil")
	}

	m.Title = "Spider-Man: Into the Spider-Verse"
	m.Released = newDate(t, "2018-12-21")

	if err := m.Validate(); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}
