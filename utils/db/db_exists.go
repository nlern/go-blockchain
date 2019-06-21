package dbutils

import "os"

const (
	dbfile = "blockchain.db"
)

// Exists checks if database exists
func Exists() bool {
	if _, err := os.Stat(dbfile); os.IsNotExist(err) {
		return false
	}
	return true
}
