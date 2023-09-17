package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
)

var (
	serverHost  = flag.String("h", "localhost", "Server Host")
	serverPort  = flag.Int("p", 25565, "Server Port")
	clientCount = flag.Int("c", 1, "Annoying Client Count")
)

func main() {
	flag.Parse()
	hostName := fmt.Sprintf("%s:%d", *serverHost, *serverPort)

	// create wait group and interrupt handler
	var wg sync.WaitGroup
	wg.Add(*clientCount)
	interruptHandler := make(chan os.Signal, 1)
	signal.Notify(interruptHandler, os.Interrupt)

	// generate clients
	fmt.Println("generating", *clientCount, "clients...")

	clients, err := GenerateClients(*clientCount, hostName, interruptHandler)
	if err != nil {
		panic(err)
	}

	// start clients
	fmt.Println("connecting to", hostName, "...")

	for _, client := range clients {
		go client.Start(&wg)
	}

	wg.Wait()
}
