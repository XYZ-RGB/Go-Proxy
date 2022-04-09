package protocol

import (
	"bytes"
	"reflect"
)

type ConnectionState int

const (
	Handshake ConnectionState = iota
	Status
	Login
	Play
)

type ConnectionDirection int

const (
	Serverbound ConnectionDirection = iota
	Clientbound
)

var Packets map[ConnectionDirection]map[ConnectionState]map[VarInt]Packet

type Packet interface {
	Write(buffer *bytes.Buffer)
	Read(session Session)
	Id() VarInt
}

func GetNewPacket(packet Packet) Packet {
	if packet == nil {
		return nil
	}

	typePacket := reflect.New(reflect.TypeOf(packet).Elem()).Interface()
	return typePacket.(Packet)
}

