package protocol

import (
	"bytes"
)

type ServerJoinGame struct {
	EntityId         Int
	Gamemode         UnsignedByte
	Dimension        Byte
	Difficulty       UnsignedByte
	MaxPlayers       UnsignedByte
	LevelType        String
	ReducedDebugInfo Boolean
}

func (s ServerJoinGame) Write(buffer *bytes.Buffer) {
	s.EntityId.Write(buffer)
	s.Gamemode.Write(buffer)
	s.Dimension.Write(buffer)
	s.Difficulty.Write(buffer)
	s.MaxPlayers.Write(buffer)
	s.LevelType.Write(buffer)
	s.ReducedDebugInfo.Write(buffer)
}

func (s *ServerJoinGame) Read(session Session) {
	//todo
}

func (s ServerJoinGame) Id() VarInt {
	return 0x01
}

type ServerSpawnPosition struct {
	Location Position
}

func (s ServerSpawnPosition) Write(buffer *bytes.Buffer) {
	s.Location.Write(buffer)
}

func (s *ServerSpawnPosition) Read(session Session) {
	//todo
}

func (s ServerSpawnPosition) Id() VarInt {
	return 0x05
}

type ServerPlayerAbilities struct {
	Invulnerable Boolean
	Flying       Boolean
	AllowFlying  Boolean
	CreativeMode Boolean
	FlyingSpeed  Float
	WalkingSpeed Float
}

func (s ServerPlayerAbilities) Write(buffer *bytes.Buffer) {
	flags := Byte(0)
	if s.Invulnerable {
		flags |= 0x01
	}
	if s.Flying {
		flags |= 0x02
	}
	if s.AllowFlying {
		flags |= 0x04
	}
	if s.CreativeMode {
		flags |= 0x08
	}
	flags.Write(buffer)
	s.FlyingSpeed.Write(buffer)
	s.WalkingSpeed.Write(buffer)
}

func (s *ServerPlayerAbilities) Read(session Session) {
	//todo
}

func (s ServerPlayerAbilities) Id() VarInt {
	return 0x39
}

type ServerPlayerPosAndLook struct {
	X     Double
	Y     Double
	Z     Double
	Yaw   Float
	Pitch Float
	Flags Byte
}

func (s ServerPlayerPosAndLook) Write(buffer *bytes.Buffer) {
	s.X.Write(buffer)
	s.Y.Write(buffer)
	s.Z.Write(buffer)
	s.Yaw.Write(buffer)
	s.Pitch.Write(buffer)
	s.Flags.Write(buffer)
}

func (s *ServerPlayerPosAndLook) Read(session Session) {
	//todo
}

func (s ServerPlayerPosAndLook) Id() VarInt {
	return 0x08
}

type ServerChatMessage struct {
	Message Msg
	Position Byte
}

func (s ServerChatMessage) Write(buffer *bytes.Buffer) {
	s.Message.Write(buffer)
	s.Position.Write(buffer)
}

func (s *ServerChatMessage) Read(session Session) {
	//todo
}

func (s ServerChatMessage) Id() VarInt {
	return 0x02
}
