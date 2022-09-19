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
	want := []byte("hello")
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Failed to deserialize bulk string, got: %s, want: %s", got, want)
	}

}
