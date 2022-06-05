package protocol

import (
	"bytes"
	"io"
)

type ClientLoginStart struct {
	Name String
}

func (c ClientLoginStart) Write(buffer *bytes.Buffer) error {
	c.Name.Write(buffer)
	return nil
}

func (c *ClientLoginStart) Read(session io.Reader) error {
	c.Name.Read(session)
	return nil
}

func (c ClientLoginStart) Id() VarInt {
	return 0x00
}

type ServerLoginSuccess struct {
	Uuid String
	Name String
}

func (s ServerLoginSuccess) Write(buffer *bytes.Buffer) error {
	s.Uuid.Write(buffer)
	s.Name.Write(buffer)
	return nil
}

func (s *ServerLoginSuccess) Read(session io.Reader) error {
	s.Uuid.Read(session)
	s.Name.Read(session)
	return nil
}

func (s ServerLoginSuccess) Id() VarInt {
	return 0x02
}