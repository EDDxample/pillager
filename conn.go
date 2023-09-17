package main

import (
	"fmt"
	"net"
	"os"

	"github.com/EDDxample/annoying_client/packet/dt"
)

type Connection struct {
	net.Conn
	addr string
	done chan os.Signal
}

func NewConnection(addr string, done chan os.Signal) *Connection {
	return &Connection{addr: addr, done: done}
}

func (c *Connection) Start() error {
	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return err
	}
	c.Conn = conn
	return nil
}

func (c *Connection) ReadPacket() ([]byte, error) {
	var length dt.VarInt
	n, err := length.ReadFrom(c)
	if err != nil {
		fmt.Printf("Client disconnected: %s (Reason: %s)\n", c.RemoteAddr(), err)
		return nil, err
	}

	buffer := make([]byte, n+int64(length))
	br, err := c.Read(buffer[n:])
	if err != nil || br == 0 {
		if err != nil {
			fmt.Printf("Client disconnected: %s (Reason: %s)\n", c.RemoteAddr(), err)
		}
		return nil, err
	}
	copy(buffer[:n+1], length.Bytes())
	return buffer, nil
}

func (c *Connection) Close() {
	c.Conn.Close()
}
