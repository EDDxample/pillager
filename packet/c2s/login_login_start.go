package c2s

import (
	"io"

	"github.com/EDDxample/pillager/packet"
	"github.com/EDDxample/pillager/packet/dt"
)

type LoginStart struct {
	Header  packet.PacketHeader
	Name    dt.String
	HasUUID dt.Boolean
	UUID    []byte
}

func (pk *LoginStart) ReadPacket(reader io.Reader) {
	pk.Header.ReadHeader(reader)
	pk.Name.ReadFrom(reader)

	pk.HasUUID.ReadFrom(reader)
	if pk.HasUUID {
		pk.UUID = make([]byte, 16)

		if _, err := reader.Read(pk.UUID); err != nil {
			panic(err)
		}
	}
}

func (pk *LoginStart) Bytes() []byte {
	pk.Header.PacketID = 0x00
	var data []byte
	data = append(data, pk.Name.Bytes()...)
	pk.HasUUID = pk.UUID != nil && len(pk.UUID) > 0
	data = append(data, pk.HasUUID.Bytes()...)
	if pk.HasUUID {
		data = append(data, pk.UUID...)
	}
	pk.Header.WriteHeader(&data)
	return data
}
