package dbutils

import "github.com/boltdb/bolt"

// Open opens a new db and returns
func Open() (*bolt.DB, error) {
	return bolt.Open(dbfile, 0600, nil)
}
