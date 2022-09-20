package miniaturedis

import (
	"reflect"
	"strings"
	"testing"
)

func TestBulkStringDeserializer(t *testing.T) {

	got, err := deserializeBulkString(strings.NewReader("$5\r\nhello\r\n"))
	if err != nil {
		t.Error("Failed to deserialize bulk string", err)
	}
	want := RespBulkString{[]byte("hello")}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Failed to deserialize bulk string, got: %s, want: %s", got, want)
	}

}

func TestBulkStringSerializer(t *testing.T) {

	tests := []struct {
		input      []byte
		serialized string
	}{
		{[]byte("hello world"), "$11\r\nhello world\r\n"},
		{[]byte(""), "$0\r\n\r\n"},
		{nil, "$-1\r\n"},
	}

	for _, test := range tests {
		bulkString := RespBulkString{[]byte(test.input)}
		got := bulkString.serialize()
		want := []byte(test.serialized)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Failed to serialize bulk string, got: %s, want: %s", got, want)
		}
	}
}

func TestArrayBulkStringDeserializer(t *testing.T) {

	got, err := deserializeArray(strings.NewReader("*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"))
	if err != nil {
		t.Error("Failed to deserialize array", err)
	}
	want := RespArray{RespBulkString{[]byte("hello")}, RespBulkString{[]byte("world")}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Failed to deserialize array, got: %s, want: %s", got, want)
	}

}
