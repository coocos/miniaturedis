package miniaturedis

import (
	"bufio"
	"errors"
	"io"
	"strconv"
)

type BulkString = []byte
type Array = []BulkString

func deserializeArray(r io.Reader) (Array, error) {

	elements := [][]byte{}

	reader := bufio.NewReader(r)
	prefix, err := reader.ReadByte()
	if prefix != '*' {
		return elements, errors.New("array should start with *")
	}
	if err != nil {
		return elements, err
	}

	// Read array size
	size := 0
	for {
		char, err := reader.ReadByte()
		if err != nil {
			return elements, err
		}
		if char == '\r' {
			break
		}
		digit, err := strconv.Atoi(string(char))
		if err != nil {
			return elements, err
		}
		size = size*10 + digit
	}
	_, err = reader.Discard(1)
	if err != nil {
		return elements, err
	}

	for len(elements) < size {
		prefix, err := reader.Peek(1)
		if err != nil {
			return elements, err
		}
		identifier := prefix[0]
		switch identifier {
		case '$':
			bulkString, err := deserializeBulkString(reader)
			if err != nil {
				return elements, err
			}
			elements = append(elements, bulkString)
		}
	}

	return elements, nil

}

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
	// Consume \n
	_, err = reader.Discard(1)
	if err != nil {
		return []byte{}, err
	}

	bulkString := make([]byte, size)
	if _, err := io.ReadFull(reader, bulkString); err != nil {
		return bulkString, err
	}

	// Consume \r\n
	_, err = reader.Discard(2)
	if err != nil {
		return []byte{}, err
	}

	return bulkString, nil
}
