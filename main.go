package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/coocos/miniaturedis/miniaturedis"
)

func main() {
	server := miniaturedis.NewServer()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		server.Stop()
	}()

	server.Start()
}
