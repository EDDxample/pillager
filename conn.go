package main

import (
	"log"
	"net"
	"time"

	"github.com/EDDxample/pillager/packet/dt"
)

type Connection struct {
	connection      net.Conn
	alive           bool
	incomingPackets chan *[]byte
	outgoingPackets chan *[]byte
}

func NewConnection(hostname string) *Connection {
	conn, err := net.Dial("tcp", hostname)
	if err != nil {
		log.Fatal(err)
	}

	instance := &Connection{
		connection:      conn,
		alive:           true,
		incomingPackets: make(chan *[]byte, 10),
		outgoingPackets: make(chan *[]byte, 10),
	}

	go instance.handleIncomingPackets()
	go instance.handleOutgoingPackets()
	return instance
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
		c.incomingPackets <- &packetBytes
	}
	c.Close()
}

func (c *Connection) handleOutgoingPackets() {
	keepAliveTicker := time.NewTicker(10 * time.Second)
	// var keepAlivePacket s2c.KeepAlive

	for c.alive {
		select {
		case packet := <-c.outgoingPackets:
			c.connection.Write(*packet)

		case t := <-keepAliveTicker.C:
			_ = t
			// if c.state == PLAY {
			// 	keepAlivePacket.KeepAliveID = dt.Long(t.UTC().UnixNano())
			// 	c.connection.Write(keepAlivePacket.Bytes())
			// }
		}
	}

	keepAliveTicker.Stop()
}

func (c *Connection) setTimeout() {
	c.connection.SetReadDeadline(time.Now().Add(10 * time.Second))
}

func (c *Connection) Close() {
	c.alive = false
	c.connection.Close()
}
