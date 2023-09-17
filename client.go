package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/EDDxample/annoying_client/packet/c2s"
	"github.com/EDDxample/annoying_client/packet/dt"
	"github.com/EDDxample/annoying_client/packet/s2c"
)

type Client struct {
	conn      *Connection
	nickname  string
	connected bool
	position  [3]dt.Double
	rotation  [2]dt.Float
	lock      sync.Mutex
}

func GenerateClients(count int, addr string, done chan os.Signal) ([]*Client, error) {
	clients := make([]*Client, count)

	for i := 0; i < count; i++ {
		conn := NewConnection(addr, done)
		clients[i] = NewClient(conn, fmt.Sprintf("bot_0%d", i))
	}

	return clients, nil
}

func NewClient(conn *Connection, nickname string) *Client {
	return &Client{conn: conn, nickname: nickname}
}

func (c *Client) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	if err := c.conn.Start(); err != nil {
		panic(err)
	}
	c.connected = true

	fmt.Println("user", c.nickname, "started.")
	c.Login()
	go c.ReadPackets()
	time.Sleep(time.Millisecond * 500)
	c.Walk()

	time.Sleep(time.Second * 1)
	c.Close()
}

func (c *Client) Login() {
	var handshake c2s.Handshake
	addr := strings.Split(c.conn.addr, ":")
	port, _ := strconv.Atoi(addr[1])
	handshake.Address = dt.String(addr[0])
	handshake.Port = dt.UShort(port)
	handshake.Protocol = dt.VarInt(762)
	handshake.NextState = 2
	c.conn.Write(handshake.Bytes())

	var loginStart c2s.LoginStart
	loginStart.Name = dt.String(c.nickname)
	c.conn.Write(loginStart.Bytes())
}

func (c *Client) Walk() {
	posPacket := c2s.SetPlayerPositionPacket{OnGround: true}
	rotPacket := c2s.SetPlayerRotationPacket{OnGround: true}
	posrotPacket := c2s.SetPlayerPositionAndRotationPacket{OnGround: true}

	for i := 0; i < 20*10; i++ {
		c.lock.Lock()

		switch i % 3 {
		case 0: // position
			c.randomWalk()
			posPacket.X = dt.Double(c.position[0])
			posPacket.Y = dt.Double(c.position[1])
			posPacket.Z = dt.Double(c.position[2])
			_, err := c.conn.Write(posPacket.Bytes())
			if err != nil {
				fmt.Println(err)
				return
			}
		case 1: // rotation
			c.randomTurn()
			rotPacket.Pitch = dt.Float(c.rotation[0])
			rotPacket.Yaw = dt.Float(c.rotation[1])
			_, err := c.conn.Write(rotPacket.Bytes())
			if err != nil {
				return
			}
		case 2: // position and rotation
			c.randomWalk()
			c.randomTurn()
			posrotPacket.X = dt.Double(c.position[0])
			posrotPacket.Y = dt.Double(c.position[1])
			posrotPacket.Z = dt.Double(c.position[2])
			posrotPacket.Yaw = dt.Float(c.rotation[0])
			posrotPacket.Pitch = dt.Float(c.rotation[1])
			_, err := c.conn.Write(posrotPacket.Bytes())
			if err != nil {
				return
			}
		}

		c.lock.Unlock()
		time.Sleep(50 * time.Millisecond)
	}
}

func (c *Client) ReadPackets() {
	for c.connected {
		p, err := c.conn.ReadPacket()
		if err != nil {
			c.connected = false
			c.Close()
			fmt.Println(err)
		}

		buffer := bytes.NewBuffer(p)
		var length dt.VarInt
		length.ReadFrom(buffer)

		var packetID dt.VarInt
		packetID.ReadFrom(buffer)

		buffer = bytes.NewBuffer(p)

		switch packetID {
		case 0x3C:
			var packet s2c.SyncPlayerPosPacket
			packet.ReadPacket(buffer)

			c.lock.Lock()
			c.position[0] = packet.X
			c.position[1] = packet.Y
			c.position[2] = packet.Z
			c.rotation[0] = packet.Yaw
			c.rotation[1] = packet.Pitch
			c.lock.Unlock()

			var response c2s.SyncPlayerPosResponsePacket
			response.TeleportID = packet.TeleportID
			c.conn.Write(response.Bytes())

		default:
			// fmt.Printf("Received %d bytes with ID: 0x%X\n", length, packetID)
		}
	}
}

func (c *Client) randomWalk() {
	const d = float64(5.7 / 20.0)

	c.position[0] += dt.Double(d * (rand.Float64()*2 - 1))
	c.position[2] += dt.Double(d * (rand.Float64()*2 - 1))
}

func (c *Client) randomTurn() {
	// yaw (XZ plane degrees)
	c.rotation[0] = dt.Float(rand.Float32() * 360)
	// pitch (90 down, -90 up)
	c.rotation[1] = 0
}

func (c *Client) Close() {
	c.connected = false
	c.conn.Close()
}
