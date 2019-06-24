package dbutils

import (
	"fmt"

	"github.com/boltdb/bolt"
)

// Open opens a new db and returns
func Open(nodeID string) (*bolt.DB, error) {
	dbFile := fmt.Sprintf(dbFile, nodeID)
	return bolt.Open(dbFile, 0600, nil)
}
