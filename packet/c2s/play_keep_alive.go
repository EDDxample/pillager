package c2s

import (
	"github.com/EDDxample/pillager/packet"
	"github.com/EDDxample/pillager/packet/dt"
)

type KeepAlivePacket struct {
	Header      packet.PacketHeader
	KeepAliveID dt.Long
}

func (pk *KeepAlivePacket) Bytes() []byte {
	pk.Header.PacketID = 0x12
	var data []byte
	data = append(data, pk.KeepAliveID.Bytes()...)
	pk.Header.WriteHeader(&data)
	return data
}
