package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/coocos/miniaturedis/internal/server"
)

func main() {
	server := server.NewServer()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		server.Stop()
	}()

	server.Start()
}
