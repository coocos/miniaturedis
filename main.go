package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os/signal"
	"sync"
	"syscall"

	"github.com/google/uuid"
)

type Server struct {
	clients  map[uuid.UUID]Client
	lock     sync.RWMutex
	requests chan Request
}

func NewServer() *Server {
	return &Server{
		make(map[uuid.UUID]Client), sync.RWMutex{}, make(chan Request),
	}
}

type Request struct {
	source uuid.UUID
	data   string
}

type Client struct {
	net.Conn
	id uuid.UUID
}

func readRequest(r io.Reader) (string, error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	if remaining := scanner.Scan(); !remaining {
		return scanner.Text(), scanner.Err()
	}
	return scanner.Text(), nil
}

func (s *Server) handleRequests() {
	for request := range s.requests {
		log.Println("Echoing", request.data)
		response := []byte(fmt.Sprintf("%s\n", request.data))
		bytes, err := s.clients[request.source].Write(response)
		if err != nil || bytes != len(response) {
			log.Println("Failed to write response to request")
		}
	}
	log.Println("Stopping request processing")
}

func (s *Server) closeConnection(client Client) {
	s.lock.Lock()
	delete(s.clients, client.id)
	s.lock.Unlock()

	err := client.Close()
	if err != nil {
		log.Println("Failed to close connection", err)
	}
	log.Println("Client dropped:", client.id)

}

func (s *Server) addConnection(conn net.Conn) {
	client := Client{conn, uuid.New()}
	log.Println("New client connection:", client.id)
	s.lock.Lock()
	s.clients[client.id] = client
	s.lock.Unlock()

	go s.handleConnection(client)
}

func (s *Server) handleConnection(client Client) {
	defer s.closeConnection(client)

	for {
		request, err := readRequest(client)
		if err != nil {
			log.Println("Invalid request from client:", err)
			return
		}
		if len(request) == 0 {
			return
		}
		s.requests <- Request{client.id, request}
		log.Println(request)
	}
}

func (s *Server) start(ctx context.Context) {

	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		log.Fatalln("Failed to start server", err)
	}

	go s.handleRequests()

	go func() {
		<-ctx.Done()
		listener.Close()
		log.Println("Stopping server")
		close(s.requests)
	}()

	log.Println("Starting server")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Failed to accept client", err)
			return
		}
		s.addConnection(conn)
	}
}

func main() {
	context, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	server := NewServer()
	server.start(context)
	stop()
}
