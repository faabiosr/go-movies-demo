package movies

import (
	"encoding/json"
	"errors"
	"time"

	v "github.com/go-ozzo/ozzo-validation/v4"
)

// Date it is a custom struct for masharlling and unmarshalling dates.
type Date struct {
	time.Time
}

// MarshalJSON marshals a Date into a "Y-m-d" format.
func (dt Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		dt.Format("2006-01-02"),
	)
}

// UnmarshalJSON unmarshals the string date into Date instance.
func (dt *Date) UnmarshalJSON(value []byte) error {
	v := string(value[1 : len(value)-1])
	if v == "" {
		return nil
	}

	t, err := time.Parse("2006-01-02", v)
	if err != nil {
		return err
	}

	dt.Time = t

	return nil
}

// Validate validates if the date is zero.
func (dt Date) Validate() error {
	if dt.IsZero() {
		return errors.New("cannot be blank")
	}

	return nil
}

// Movie represents a movie in the system.
type Movie struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Released Date   `json:"released"`
}

// Movies represents a collection of movie.
type Movies []Movie

// Validate validates the movie structure.
func (m Movie) Validate() error {
	const min = 3

	return v.ValidateStruct(&m,
		v.Field(&m.Title, v.Required, v.Length(min, 0)),
		v.Field(&m.Released, v.Required),
	)
}
