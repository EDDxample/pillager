package c2s

import (
	"github.com/EDDxample/annoying_client/packet"
	"github.com/EDDxample/annoying_client/packet/dt"
)

type SetPlayerRotationPacket struct {
	Header   packet.PacketHeader
	Yaw      dt.Float
	Pitch    dt.Float
	OnGround dt.Boolean
}

func (pk *SetPlayerRotationPacket) Bytes() []byte {
	pk.Header.PacketID = 0x16
	var data []byte
	data = append(data, pk.Yaw.Bytes()...)
	data = append(data, pk.Pitch.Bytes()...)
	data = append(data, pk.OnGround.Bytes()...)
	pk.Header.WriteHeader(&data)
	return data
}
