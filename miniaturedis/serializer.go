package miniaturedis

import (
	"bytes"
	"fmt"
)

func serializeBulkString(bulkString BulkString) []byte {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("$%d\r\n", len(bulkString)))
	buffer.Write(bulkString)
	buffer.WriteString("\r\n")
	return buffer.Bytes()
}
