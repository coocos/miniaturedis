package miniaturedis

import (
	"errors"
	"fmt"
	"log"
)

type Database interface {
	get(key string) (string, error)
	execute(request RespArray) RespType
}

type DefaultDatabase struct {
	data map[string]string
}

func NewDatabase() *DefaultDatabase {
	return &DefaultDatabase{
		make(map[string]string),
	}
}

var ErrKeyDoesNotExist = errors.New("key does not exist")

func (d *DefaultDatabase) get(key string) (string, error) {
	value, ok := d.data[key]
	if !ok {
		return "", ErrKeyDoesNotExist
	}
	return value, nil
}

func (d *DefaultDatabase) execute(request RespArray) RespType {
	command := request[0].data
	switch string(command) {
	case "GET":
		{
			value, err := d.get(string(request[1].data))
			if err != nil && errors.Is(err, ErrKeyDoesNotExist) {
				return RespBulkString{}
			} else if err != nil {
				log.Println("Failed to execute GET", err)
				return RespError{message: "something went wrong"}
			} else {
				return RespBulkString{[]byte(value)}
			}
		}
	}
	return RespError{message: fmt.Sprintf("unknown command %s", command)}
}
