package c2s

import (
	"github.com/EDDxample/annoying_client/packet"
	"github.com/EDDxample/annoying_client/packet/dt"
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
