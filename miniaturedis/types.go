package miniaturedis

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
)

type RespType interface {
	serialize() []byte
}

type RespArray = []RespBulkString

type RespBulkString struct {
	data []byte
}

func (s RespBulkString) serialize() []byte {
	if s.data == nil {
		return []byte("$-1\r\n")
	}
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("$%d\r\n", len(s.data)))
	buffer.Write(s.data)
	buffer.WriteString("\r\n")
	return buffer.Bytes()
}

type RespError struct {
	message string
}

func (s RespError) serialize() []byte {
	return []byte(fmt.Sprintf("-ERR %s\r\n", s.message))
}

func deserializeArray(r io.Reader) (RespArray, error) {

	array := RespArray{}

	reader := bufio.NewReader(r)
	prefix, err := reader.ReadByte()
	if err != nil {
		return array, err
	}
	if prefix != '*' {
		return array, fmt.Errorf("array should start with *, not %v", prefix)
	}

	// Read array size
	size := 0
	for {
		char, err := reader.ReadByte()
		if err != nil {
			return array, err
		}
		if char == '\r' {
			break
		}
		digit, err := strconv.Atoi(string(char))
		if err != nil {
			return array, err
		}
		size = size*10 + digit
	}
	_, err = reader.Discard(1)
	if err != nil {
		return array, err
	}

	for len(array) < size {
		prefix, err := reader.Peek(1)
		if err != nil {
			return array, err
		}
		identifier := prefix[0]
		switch identifier {
		case '$':
			bulkString, err := deserializeBulkString(reader)
			if err != nil {
				return array, err
			}
			array = append(array, bulkString)
		}
	}

	return array, nil

}

func deserializeBulkString(r io.Reader) (RespBulkString, error) {

	bulkString := RespBulkString{}
	reader := bufio.NewReader(r)

	prefix, err := reader.ReadByte()
	if prefix != '$' {
		return bulkString, errors.New("bulk string should start with $")
	}
	if err != nil {
		return bulkString, err
	}

	// Read string size
	size := 0
	for {
		char, err := reader.ReadByte()
		if err != nil {
			return bulkString, err
		}
		if char == '\r' {
			break
		}
		digit, err := strconv.Atoi(string(char))
		if err != nil {
			return bulkString, err
		}
		size = size*10 + digit
	}
	// Consume \n
	_, err = reader.Discard(1)
	if err != nil {
		return bulkString, err
	}

	bulkString.data = make([]byte, size)
	if _, err := io.ReadFull(reader, bulkString.data); err != nil {
		return bulkString, err
	}

	// Consume \r\n
	_, err = reader.Discard(2)
	if err != nil {
		return bulkString, err
	}

	return bulkString, nil
}
