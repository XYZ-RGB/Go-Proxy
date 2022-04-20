package protocol

import (
	"bytes"
)

type ClientLoginStart struct {
	Name String
}

func (c ClientLoginStart) Write(buffer *bytes.Buffer) {
	c.Name.Write(buffer)
}

func (c *ClientLoginStart) Read(session Session) {
	c.Name.Read(session)
}

func (c ClientLoginStart) Id() VarInt {
	return 0x00
}

type ServerLoginSuccess struct {
	Uuid String
	Name String
}

func (s ServerLoginSuccess) Write(buffer *bytes.Buffer) {
	s.Uuid.Write(buffer)
	s.Name.Write(buffer)
}

func (s *ServerLoginSuccess) Read(session Session) {
	s.Uuid.Read(session)
	s.Name.Read(session)
}

func (s ServerLoginSuccess) Id() VarInt {
	return 0x02
}