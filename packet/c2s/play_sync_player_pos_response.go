package c2s

import (
	"github.com/EDDxample/pillager/packet"
	"github.com/EDDxample/pillager/packet/dt"
)

type SyncPlayerPosResponsePacket struct {
	Header     packet.PacketHeader
	TeleportID dt.VarInt
}

func (pk *SyncPlayerPosResponsePacket) Bytes() []byte {
	pk.Header.PacketID = 0x00
	data := pk.TeleportID.Bytes()
	pk.Header.WriteHeader(&data)
	return data
}
