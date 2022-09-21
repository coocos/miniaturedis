package server

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"testing"
	"time"
)

func NewTestClient() (net.Conn, error) {
	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		return nil, err
	}
	err = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	return conn, err
}

func TestGetRequest(t *testing.T) {

	database := NewDatabase()
	database.data["example-key"] = "example-value"
	server := NewServerWithDatabase(database)
	go server.Start()
	defer server.Stop()

	client, err := NewTestClient()
	if err != nil {
		t.Error("Failed to connect to server")
	}
	defer client.Close()

	request := []byte("*2\r\n$3\r\nGET\r\n$3\r\nabc\r\n")
	if _, err := client.Write(request); err != nil {
		t.Fatal("Failed to send request to server", err)
	}

	want := []byte("$-1\r\n")
	got, err := bufio.NewReader(client).ReadBytes('\n')
	if err != nil {
		t.Fatal("Failed to receive response:", err)
	}
	if !bytes.Equal(got, want) {
		t.Errorf("Expected %s, got %s", want, got)
	}

	request = []byte("*2\r\n$3\r\nGET\r\n$11\r\nexample-key\r\n")
	if _, err := client.Write(request); err != nil {
		t.Fatal("Failed to send request to server", err)
	}

	want = []byte("$13\r\nexample-value\r\n")
	got = make([]byte, len(want))
	_, err = io.ReadFull(client, got)
	if err != nil {
		t.Fatal("Failed to receive response:", err)
	}
	if !bytes.Equal(got, want) {
		t.Errorf("Expected %s, got %s", want, got)
	}

}

func TestSetRequest(t *testing.T) {

	database := NewDatabase()
	server := NewServerWithDatabase(database)
	go server.Start()
	defer server.Stop()

	client, err := NewTestClient()
	if err != nil {
		t.Error("Failed to connect to server")
	}
	defer client.Close()

	request := []byte("*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n")
	if _, err := client.Write(request); err != nil {
		t.Fatal("Failed to send request to server", err)
	}

	want := []byte("+OK\r\n")
	got, err := bufio.NewReader(client).ReadBytes('\n')
	if err != nil {
		t.Fatal("Failed to receive response:", err)
	}
	if !bytes.Equal(got, want) {
		t.Errorf("Expected %s, got %s", want, got)
	}

	if database.data["key"] != "value" {
		t.Error("Key was not set")
	}

}
