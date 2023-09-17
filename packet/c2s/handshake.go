package c2s

import (
	"io"

	"github.com/EDDxample/annoying_client/packet"
	"github.com/EDDxample/annoying_client/packet/dt"
)

type Handshake struct {
	Header    packet.PacketHeader
	Protocol  dt.VarInt
	Address   dt.String
	Port      dt.UShort
	NextState dt.VarInt
}

func (pk *Handshake) ReadPacket(reader io.Reader) {
	pk.Header.ReadHeader(reader)
	pk.Protocol.ReadFrom(reader)
	pk.Address.ReadFrom(reader)
	pk.Port.ReadFrom(reader)
	pk.NextState.ReadFrom(reader)
}

func (pk *Handshake) Bytes() []byte {
	pk.Header.PacketID = 0x00
	var data []byte
	data = append(data, pk.Protocol.Bytes()...)
	data = append(data, pk.Address.Bytes()...)
	data = append(data, pk.Port.Bytes()...)
	data = append(data, pk.NextState.Bytes()...)
	pk.Header.WriteHeader(&data)
	return data
}
