package s2c

import (
	"io"

	"github.com/EDDxample/pillager/packet"
	"github.com/EDDxample/pillager/packet/dt"
)

type KeepAlivePacket struct {
	Header      packet.PacketHeader
	KeepAliveID dt.Long
}

func (pk *KeepAlivePacket) ReadPacket(reader io.Reader) {
	pk.Header.ReadHeader(reader)
	pk.KeepAliveID.ReadFrom(reader)
}
