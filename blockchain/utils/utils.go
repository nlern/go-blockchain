/*
Package utils provide helper utilities
*/
package utils

import (
	"bytes"
	"encoding/binary"
	"log"
)

/*
IntToHex method converts a int64 number to byte array
*/
func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)

	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}
