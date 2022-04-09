package main

import (
	"Proxy/protocol"
	"net"
	"time"
)

func main() {
	handshakeMap := map[protocol.VarInt]protocol.Packet{
		protocol.HandshakePacket{}.Id(): &protocol.HandshakePacket{},
	}
	loginMap := map[protocol.VarInt]protocol.Packet{
		protocol.ClientLoginStart{}.Id(): &protocol.ClientLoginStart{},
	}
	statusMap := map[protocol.VarInt]protocol.Packet{
		protocol.ClientStatusRequest{}.Id(): &protocol.ClientStatusRequest{},
	}

	playMap := map[protocol.VarInt]protocol.Packet{
		protocol.ClientChatMessage{}.Id(): &protocol.ClientChatMessage{},
		protocol.KeepAlive{}.Id():         &protocol.KeepAlive{},
	}

	serverBoundMap := map[protocol.ConnectionState]map[protocol.VarInt]protocol.Packet{
		protocol.Handshake: handshakeMap,
		protocol.Login:     loginMap,
		protocol.Status:    statusMap,
		protocol.Play:      playMap,
	}

	protocol.Packets = map[protocol.ConnectionDirection]map[protocol.ConnectionState]map[protocol.VarInt]protocol.Packet{
		protocol.Serverbound: serverBoundMap,
		//todo protocol.Clientbound:
	}

	listen, err := net.Listen("tcp", ":25565")
	if err != nil {
		return
	}

	go func() {
		for range time.Tick(time.Second * 5) {
			for _, session := range protocol.Sessions {
				session.SendPacket(&protocol.KeepAlive{KeepAliveId: protocol.VarInt(time.Now().Unix())})
			}
		}
	}()

	for {
		conn, err := listen.Accept()
		if err != nil {
			return
		}

		go protocol.HandleConnection(conn)
	}
}
