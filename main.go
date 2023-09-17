package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	host  = flag.String("h", ":25565", "Server Hostname")
	count = flag.Int("c", 1, "Pillager Count")
)

func main() {
	flag.Parse()
	log.Printf("Raiding server %s with %d pillager(s)...\n", *host, *count)

	clients := CreateClients(*host, *count)
	<-handleInterrupt()
	shutdown(clients)

	log.Println("done")
}

func handleInterrupt() chan os.Signal {
	endSignal := make(chan os.Signal, 2)
	signal.Notify(endSignal, syscall.SIGINT, syscall.SIGTERM)
	return endSignal
}

func shutdown(clients []*Client) {
	for i := 0; i < len(clients); i++ {
		clients[i].Close()
		clients[i] = nil
	}
}
