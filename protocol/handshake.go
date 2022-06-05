package protocol

import (
	"bytes"
	"io"
)

type HandshakePacket struct {
	ProtocolVersion VarInt
	ServerAddress   String
	ServerPort      UnsignedShort
	NextState       VarInt
}

func (h HandshakePacket) Write(buffer *bytes.Buffer) error {
	h.ProtocolVersion.Write(buffer)
	h.ServerAddress.Write(buffer)
	h.ServerPort.Write(buffer)
	h.NextState.Write(buffer)
	return nil
}

func (h *HandshakePacket) Read(session io.Reader) error {
	h.ProtocolVersion.Read(session)
	h.ServerAddress.Read(session)
	h.ServerPort.Read(session)
	h.NextState.Read(session)
	return nil
}

func (h HandshakePacket) Id() VarInt {
	return 0x00
}
