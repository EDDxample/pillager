package packet

import (
	"io"

	"github.com/EDDxample/annoying_client/packet/dt"
)

type PacketHeader struct {
	Length   dt.VarInt
	PacketID dt.VarInt
}

func (pb *PacketHeader) ReadHeader(reader io.Reader) {
	pb.Length.ReadFrom(reader)
	pb.PacketID.ReadFrom(reader)
}

func (pb *PacketHeader) WriteHeader(buffer *[]byte) {
	*buffer = append(pb.PacketID.Bytes(), *buffer...)
	pb.Length = dt.VarInt(len(*buffer))
	*buffer = append(pb.Length.Bytes(), *buffer...)
}
