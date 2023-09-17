package s2c

import (
	"io"

	"github.com/EDDxample/annoying_client/packet"
	"github.com/EDDxample/annoying_client/packet/dt"
)

type SyncPlayerPosPacket struct {
	Header     packet.PacketHeader
	X          dt.Double
	Y          dt.Double
	Z          dt.Double
	Yaw        dt.Float
	Pitch      dt.Float
	Flags      dt.Byte
	TeleportID dt.VarInt
}

func (pk *SyncPlayerPosPacket) ReadPacket(reader io.Reader) {
	pk.Header.ReadHeader(reader)
	pk.X.ReadFrom(reader)
	pk.Y.ReadFrom(reader)
	pk.Z.ReadFrom(reader)
	pk.Yaw.ReadFrom(reader)
	pk.Pitch.ReadFrom(reader)
	pk.Flags.ReadFrom(reader)
	pk.TeleportID.ReadFrom(reader)
}
