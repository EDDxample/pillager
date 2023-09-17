package c2s

import (
	"github.com/EDDxample/pillager/packet"
	"github.com/EDDxample/pillager/packet/dt"
)

type SetPlayerPositionAndRotationPacket struct {
	Header   packet.PacketHeader
	X        dt.Double
	Y        dt.Double
	Z        dt.Double
	Yaw      dt.Float
	Pitch    dt.Float
	OnGround dt.Boolean
}

func (pk *SetPlayerPositionAndRotationPacket) Bytes() []byte {
	pk.Header.PacketID = 0x15
	var data []byte
	data = append(data, pk.X.Bytes()...)
	data = append(data, pk.Y.Bytes()...)
	data = append(data, pk.Z.Bytes()...)
	data = append(data, pk.Yaw.Bytes()...)
	data = append(data, pk.Pitch.Bytes()...)
	data = append(data, pk.OnGround.Bytes()...)
	pk.Header.WriteHeader(&data)
	return data
}
