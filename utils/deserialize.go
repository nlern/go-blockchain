package utils

import (
	"bytes"
	"encoding/gob"
)

// Deserialize deserializes data into targetPointer and returns any error
func Deserialize(registerValue interface{}, data []byte, targetPointer interface{}) error {
	if registerValue != nil {
		gob.Register(registerValue)
	}
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(targetPointer)

	return err
}
