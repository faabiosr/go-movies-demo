package movies

import (
	"encoding/json"
	"fmt"

	"github.com/gofrs/uuid/v5"
	bolt "go.etcd.io/bbolt"

	"github.com/faabiosr/go-movies-demo/pkg/errors"
)

// Datasource manages the movies data.
type Datasource struct {
	db *bolt.DB
}

var bucketName = []byte("movies")

// NewDatasource creates a new instance of movie datasource.
func NewDatasource(db *bolt.DB) *Datasource {
	return &Datasource{db: db}
}

// Store stores the movie content into the database, if some error occurs, an
// error will be returned.
func (ds *Datasource) Store(m Movie) (Movie, error) {
	if m.ID == "" {
		m.ID = uuid.Must(uuid.NewV4()).String()
	}

	err := ds.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(bucketName)
		if err != nil {
			return err
		}

		entry, err := json.Marshal(m)
		if err != nil {
			return err
		}

		return bucket.Put([]byte(m.ID), entry)
	})

	return m, err
}

// Fetch fetches the movie by id, if not exist an error will be returned.
func (ds *Datasource) Fetch(id string) (Movie, error) {
	var entry []byte
	var m Movie

	err := ds.db.View(func(tx *bolt.Tx) error {
		if bucket := tx.Bucket(bucketName); bucket != nil {
			entry = bucket.Get([]byte(id))
			return nil
		}

		return errors.New("movie was not found")
	})
	if err != nil {
		return m, err
	}

	if len(entry) == 0 {
		return m, fmt.Errorf("movie was not found")
	}

	return m, json.Unmarshal(entry, &m)
}

// FetchAll retrieves all the movies stored.
func (ds *Datasource) FetchAll() (Movies, error) {
	movies := Movies{}

	err := ds.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		if bucket == nil {
			return nil
		}

		return bucket.ForEach(func(_, v []byte) error {
			var m Movie

			if err := json.Unmarshal(v, &m); err != nil {
				return err
			}

			movies = append(movies, m)

			return nil
		})
	})

	return movies, err
}

// Remove removes the movie from database, if some error occurs, an error will
// be returned.
func (ds *Datasource) Remove(id string) error {
	if id == "" {
		return errors.New("empty id is not allowed")
	}

	return ds.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		if bucket == nil {
			return errors.New("movie was not found")
		}

		return bucket.Delete([]byte(id))
	})
}
