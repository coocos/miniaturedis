package miniaturedis

import (
	"bufio"
	"errors"
	"io"
	"strconv"
)

type BulkString = []byte

func deserializeBulkString(r io.Reader) (BulkString, error) {

	reader := bufio.NewReader(r)

	prefix, err := reader.ReadByte()
	if prefix != '$' {
		return []byte{}, errors.New("bulk string should start with $")
	}
	if err != nil {
		return []byte{}, err
	}

	// Read string size
	size := 0
	for {
		char, err := reader.ReadByte()
		if err != nil {
			return []byte{}, err
		}
		if char == '\r' {
			break
		}
		digit, err := strconv.Atoi(string(char))
		if err != nil {
			return []byte{}, err
		}
		size = size*10 + digit
	}
	_, err = reader.Discard(1)
	if err != nil {
		return []byte{}, err
	}

	bulkString := make([]byte, size)
	if _, err := io.ReadFull(reader, bulkString); err != nil {
		return bulkString, err
	}

	return bulkString, nil
}
