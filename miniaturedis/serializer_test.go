package miniaturedis

import (
	"reflect"
	"testing"
)

func TestBulkStringSerializer(t *testing.T) {

	tests := []struct {
		input      string
		serialized string
	}{
		{"hello world", "$11\r\nhello world\r\n"},
		{"", "$0\r\n\r\n"},
	}

	for _, test := range tests {
		got := serializeBulkString([]byte(test.input))
		want := []byte(test.serialized)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Failed to serialize bulk string, got: %s, want: %s", got, want)
		}
	}
}
