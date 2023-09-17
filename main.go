package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	host  = flag.String("h", ":25565", "Server Hostname")
	count = flag.Int("c", 1, "Pillager Count")
)

func main() {
	flag.Parse()
	log.Printf("Raiding server %s with %d pillager(s)...\n", *host, *count)

	var wg sync.WaitGroup
	wg.Add(*count)
	clients := CreateClients(*host, *count, &wg)

	wg.Add(*count)
	<-handleInterrupt(&wg)
	shutdown(clients)

	log.Println("done")
}

func handleInterrupt(wg *sync.WaitGroup) chan os.Signal {
	endSignal := make(chan os.Signal, 2)
	signal.Notify(endSignal, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		wg.Wait()
		endSignal <- syscall.SIGINT
	}()

	return endSignal
}

func shutdown(clients []*Client) {
	for i := 0; i < len(clients); i++ {
		clients[i].Close()
		clients[i] = nil
	}
}
