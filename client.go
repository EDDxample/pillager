package main

import (
	"fmt"
	"log"
	"sync"
)

type Client struct {
	nickname string
	conn     *Connection
}

func CreateClients(hostname string, count int) []*Client {
	var wg sync.WaitGroup
	wg.Add(count)

	clients := make([]*Client, count)
	for i := 0; i < count; i++ {
		go func(id int) {
			log.Printf("Starting pillager Nº%d.\n", id)
			nickname := fmt.Sprintf("pillager_%d", id)
			clients[id] = NewClient(nickname, hostname)
			wg.Done()
		}(i)
	}
	wg.Wait()

	return clients
}

func NewClient(nickname, hostname string) *Client {
	conn := NewConnection(hostname)
	return &Client{
		nickname: "pillager",
		conn:     conn,
	}
}

func Start() {
}

func Login() {}

func Play() {}

func (c *Client) Close() {
	c.conn.Close()
}
