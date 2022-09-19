package main

import (
	"bufio"
	"bytes"
	"log"
	"net"
	"testing"
)

func NewTestClient() (net.Conn, error) {
	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func TestEcho(t *testing.T) {

	server := NewServer()
	go server.start()
	defer server.stop()

	client, err := NewTestClient()
	if err != nil {
		t.Error("Failed to connect to server")
	}
	defer client.Close()

	log.Println("Writing to server")
	message := []byte("hello\n")
	if _, err := client.Write(message); err != nil {
		t.Error("Failed to send message to server", err)
	}

	log.Println("Waiting for response from server")
	response, err := bufio.NewReader(client).ReadBytes('\n')
	if err != nil || !bytes.Equal(message, response) {
		t.Error("Failed to receive echo response")
	}

}
