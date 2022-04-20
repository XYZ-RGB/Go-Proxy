package protocol

import (
	"bytes"
)

type HandshakePacket struct {
	ProtocolVersion VarInt
	ServerAddress   String
	ServerPort      UnsignedShort
	NextState       VarInt
}

func (h HandshakePacket) Write(buffer *bytes.Buffer) {
	h.ProtocolVersion.Write(buffer)
	h.ServerAddress.Write(buffer)
	h.ServerPort.Write(buffer)
	h.NextState.Write(buffer)
}

func (h *HandshakePacket) Read(session Session) {
	h.ProtocolVersion.Read(session)
	h.ServerAddress.Read(session)
	h.ServerPort.Read(session)
	h.NextState.Read(session)
}

func (h HandshakePacket) Id() VarInt {
	return 0x00
}
