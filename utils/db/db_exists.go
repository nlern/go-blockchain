package dbutils

import (
	"fmt"
	"os"
)

const (
	dbFile = "blockchain_%s.db"
)

// Exists checks if database exists
func Exists(nodeID string) bool {
	dbFile := fmt.Sprintf(dbFile, nodeID)
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}
