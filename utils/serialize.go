package utils

import (
	"bytes"
	"encoding/gob"
)

// Serialize serializes data and returns it with any error
func Serialize(registerValue interface{}, data interface{}) ([]byte, error) {
	var result bytes.Buffer

	if registerValue != nil {
		gob.Register(registerValue)
	}

	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(data)

	return result.Bytes(), err
}
