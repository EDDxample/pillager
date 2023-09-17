package c2s

import (
	"github.com/EDDxample/annoying_client/packet"
	"github.com/EDDxample/annoying_client/packet/dt"
)

type SetPlayerPositionPacket struct {
	Header   packet.PacketHeader
	X        dt.Double
	Y        dt.Double
	Z        dt.Double
	OnGround dt.Boolean
}

func (pk *SetPlayerPositionPacket) Bytes() []byte {
	pk.Header.PacketID = 0x14
	var data []byte
	data = append(data, pk.X.Bytes()...)
	data = append(data, pk.Y.Bytes()...)
	data = append(data, pk.Z.Bytes()...)
	data = append(data, pk.OnGround.Bytes()...)
	pk.Header.WriteHeader(&data)
	return data
}
