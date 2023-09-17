package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/EDDxample/pillager/packet/c2s"
	"github.com/EDDxample/pillager/packet/dt"
	"github.com/EDDxample/pillager/packet/s2c"
)

const tickDuration = 50 * time.Millisecond
const sprintSpeed = float64(5.7 / 20.0)

type Client struct {
	nickname   string
	conn       *Connection
	wg         *sync.WaitGroup
	x, y, z    float64
	yaw, pitch float32
	spawned    bool
}

func CreateClients(hostname string, count int, wg *sync.WaitGroup) []*Client {
	clients := make([]*Client, count)
	for i := 0; i < count; i++ {
		go func(id int) {
			log.Printf("Starting pillager Nº%d.\n", id)
			nickname := fmt.Sprintf("pillager_%d", id)
			clients[id] = NewClient(nickname, hostname, wg)
			wg.Done()
		}(i)
	}
	wg.Wait()

	for _, client := range clients {
		go client.start()
	}

	return clients
}

func NewClient(nickname, hostname string, wg *sync.WaitGroup) *Client {
	conn := NewConnection(hostname)
	return &Client{
		nickname: nickname,
		conn:     conn,
		wg:       wg,
	}
}

func (c *Client) start() {
	defer c.wg.Done()
	c.login()
	time.Sleep(10 * tickDuration)
	c.play()
}

func (c *Client) login() {
	host, port := c.conn.GetHostPort()

	handshake := c2s.Handshake{
		Address:   dt.String(host),
		Port:      dt.UShort(port),
		Protocol:  762, // 1.19.4
		NextState: 2,   // login
	}
	c.conn.outgoingPackets <- handshake.Bytes()

	loginStart := c2s.LoginStart{Name: dt.String(c.nickname)}
	c.conn.outgoingPackets <- loginStart.Bytes()
}

func (c *Client) play() {
	posPacket := c2s.SetPlayerPositionPacket{OnGround: true}
	rotPacket := c2s.SetPlayerRotationPacket{OnGround: true}
	posrotPacket := c2s.SetPlayerPositionAndRotationPacket{OnGround: true}
	for i := 0; c.conn.alive; i++ {
		c.handleIncomingPackets()
		if c.spawned {
			c.randomWalk()
			switch i % 3 {
			case 0: // position
				posPacket.X = dt.Double(c.x)
				posPacket.Y = dt.Double(c.y)
				posPacket.Z = dt.Double(c.z)
				c.conn.outgoingPackets <- posPacket.Bytes()

			case 1: // rotation
				rotPacket.Yaw = dt.Float(c.yaw)
				rotPacket.Pitch = dt.Float(c.pitch)
				c.conn.outgoingPackets <- rotPacket.Bytes()

			case 2: // position and rotation
				posrotPacket.X = dt.Double(c.x)
				posrotPacket.Y = dt.Double(c.y)
				posrotPacket.Z = dt.Double(c.z)
				posrotPacket.Yaw = dt.Float(c.yaw)
				posrotPacket.Pitch = dt.Float(c.pitch)
				c.conn.outgoingPackets <- posrotPacket.Bytes()
			}
		}
		time.Sleep(tickDuration)
	}
}

func (c *Client) randomWalk() {
	c.x += sprintSpeed * (rand.Float64()*2 - 1)
	c.z += sprintSpeed * (rand.Float64()*2 - 1)
	c.yaw = rand.Float32() * 360
	c.pitch = (rand.Float32()*2 - 1) * 90
}

func (c *Client) handleIncomingPackets() {
	for {
		select {
		case packet := <-c.conn.incomingPackets:
			buffer := bytes.NewBuffer(packet)
			var length dt.VarInt
			length.ReadFrom(buffer)

			var packetID dt.VarInt
			packetID.ReadFrom(buffer)
			buffer = bytes.NewBuffer(packet)

			switch packetID {
			case 0x3C:
				c.handleTeleportPacket(buffer)
			case 0x23:
				c.handleKeepAlivePacket(buffer)
			}
		default:
			return
		}
	}
}

func (c *Client) handleTeleportPacket(reader io.Reader) {
	var packet s2c.SyncPlayerPosPacket
	packet.ReadPacket(reader)

	c.x = float64(packet.X)
	c.y = float64(packet.Y)
	c.z = float64(packet.Z)
	c.yaw = float32(packet.Yaw)
	c.pitch = float32(packet.Pitch)
	c.spawned = true

	response := c2s.SyncPlayerPosResponsePacket{TeleportID: packet.TeleportID}
	c.conn.outgoingPackets <- response.Bytes()
}

func (c *Client) handleKeepAlivePacket(reader io.Reader) {
	var packet s2c.KeepAlivePacket
	packet.ReadPacket(reader)
	response := c2s.KeepAlivePacket{KeepAliveID: packet.KeepAliveID}
	c.conn.outgoingPackets <- response.Bytes()
}

func (c *Client) Close() {
	c.conn.Close()
}
