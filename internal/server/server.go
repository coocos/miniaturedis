package server

import (
	"errors"
	"io"
	"log"
	"net"
	"sync"

	"github.com/google/uuid"
)

type Server struct {
	clients  map[uuid.UUID]Client
	lock     sync.RWMutex
	requests chan Request
	shutdown chan any
	listener net.Listener
	database Database
}

func NewServer() *Server {
	return &Server{
		clients:  make(map[uuid.UUID]Client),
		requests: make(chan Request),
		shutdown: make(chan any),
		database: NewDatabase(),
	}
}

func NewServerWithDatabase(database Database) *Server {
	return &Server{
		clients:  make(map[uuid.UUID]Client),
		requests: make(chan Request),
		shutdown: make(chan any),
		database: database,
	}
}

type Request struct {
	source uuid.UUID
	data   RespArray
}

type Client struct {
	net.Conn
	id        uuid.UUID
	responses chan RespType
}

func readRequest(r io.Reader) (RespArray, error) {
	array, err := deserializeArray(r)
	if err != nil {
		return RespArray{}, err
	}
	return array, nil
}

func (s *Server) handleRequests() {
	for request := range s.requests {
		log.Printf("Client %s sent %s\n", request.source, request.data)
		response := s.database.execute(request.data)
		log.Printf("Response: %s", response)

		s.lock.RLock()
		s.clients[request.source].responses <- response
		s.lock.RUnlock()
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
	client := Client{conn, uuid.New(), make(chan RespType)}
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
			if errors.Is(err, io.EOF) {
				return
			} else {
				log.Println("Invalid request from client:", err)
				return
			}
		}
		s.requests <- Request{client.id, request}
		response := <-client.responses
		bytes, err := client.Write(response.serialize())
		if err != nil || bytes != len(response.serialize()) {
			log.Printf("Failed to send response to client %s: %v", client.id, err)
			return
		}
	}
}

func (s *Server) Stop() {
	log.Println("Stopping server")
	close(s.shutdown)
	s.listener.Close()
	close(s.requests)
}

// Start starts the server
func (s *Server) Start() {

	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		log.Fatalln("Failed to start server", err)
	}
	s.listener = listener

	go s.handleRequests()

	log.Println("Starting server")
	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-s.shutdown:
				return
			default:
				log.Println("Failed to accept client", err)
			}
		} else {
			s.addConnection(conn)
		}
	}
}
