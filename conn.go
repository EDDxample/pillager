package main

import (
	"log"
	"net"
	"strconv"
	"time"

	"github.com/EDDxample/pillager/packet/dt"
)

type Connection struct {
	connection      net.Conn
	alive           bool
	incomingPackets chan []byte
	outgoingPackets chan []byte
}

func NewConnection(hostname string) *Connection {
	hostname = ensurePort(hostname, "25565")
	conn, err := net.Dial("tcp", hostname)
	if err != nil {
		log.Fatal(err)
	}

	instance := &Connection{
		connection:      conn,
		alive:           true,
		incomingPackets: make(chan []byte, 50),
		outgoingPackets: make(chan []byte, 10),
	}

	go instance.handleIncomingPackets()
	go instance.handleOutgoingPackets()
	return instance
}

func ensurePort(host, defaultPort string) string {
	if _, _, err := net.SplitHostPort(host); err != nil {
		if addrErr, ok := err.(*net.AddrError); ok && addrErr.Err == "missing port in address" {
			return net.JoinHostPort(host, defaultPort)
		}
	}
	return host
}

func (c *Connection) GetHostPort() (string, int) {
	host, portStr, err := net.SplitHostPort(c.connection.RemoteAddr().String())
	if err != nil {
		log.Fatal(err)
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatal(err)
	}

	return host, port
}

func (c *Connection) handleIncomingPackets() {
	for c.alive {
		c.setTimeout()

		var packetLength dt.VarInt
		n, err := packetLength.ReadFrom(c.connection)
		if err != nil || n == 0 || packetLength == 0 {
			if err != nil && c.alive {
				log.Printf("Disconnected: %s (Reason: %s)\n", c.connection.RemoteAddr(), err)
			}
			break
		}

		packetBytes := make([]byte, n+int64(packetLength))
		br, err := c.connection.Read(packetBytes[n:])
		if err != nil || br == 0 {
			if err != nil && c.alive {
				log.Printf("Disconnected: %s (Reason: %s)\n", c.connection.RemoteAddr(), err)
			}
			break
		}

		copy(packetBytes[:n+1], packetLength.Bytes())
		c.incomingPackets <- packetBytes
	}
	c.Close()
}

func (c *Connection) handleOutgoingPackets() {
	for c.alive {
		select {
		case packet := <-c.outgoingPackets:
			c.connection.Write(packet)
		}
	}
}

func (c *Connection) setTimeout() {
	c.connection.SetReadDeadline(time.Now().Add(10 * time.Second))
}

func (c *Connection) Close() {
	c.alive = false
	c.connection.Close()
}
