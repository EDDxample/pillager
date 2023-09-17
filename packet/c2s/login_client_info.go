package c2s

import (
	"github.com/EDDxample/pillager/packet"
	"github.com/EDDxample/pillager/packet/dt"
)

type ClientInfoPacket struct {
	Header              packet.PacketHeader
	Locale              dt.String
	ViewDistance        dt.Byte
	ChatMode            dt.VarInt
	ChatColors          dt.Boolean
	DisplayedSkinParts  dt.UByte
	MainHand            dt.VarInt
	EnableTextFiltering dt.Boolean
	AllowServerListings dt.Boolean
}

func (pk *ClientInfoPacket) Bytes() []byte {
	pk.Header.PacketID = 0x08
	var data []byte
	data = append(data, pk.Locale.Bytes()...)
	data = append(data, pk.ViewDistance.Bytes()...)
	data = append(data, pk.ChatMode.Bytes()...)
	data = append(data, pk.ChatColors.Bytes()...)
	data = append(data, pk.DisplayedSkinParts.Bytes()...)
	data = append(data, pk.MainHand.Bytes()...)
	data = append(data, pk.EnableTextFiltering.Bytes()...)
	data = append(data, pk.AllowServerListings.Bytes()...)
	pk.Header.WriteHeader(&data)
	return data
}
