package protocol

import (
	"bytes"
	"github.com/google/uuid"
	"io"
	"net"
)

var Sessions = make(map[string]*Session)

type Session struct {
	net.Conn
	io.Reader
	io.Writer
	state     ConnectionState
	direction ConnectionDirection
	uuid      uuid.UUID
}

func (s Session) close() {
	Sessions[s.uuid.String()] = nil
	s.Conn.Close()
}

func (s Session) SendPacket(packet Packet) {
	tempBuf := new(bytes.Buffer)
	packet.Id().Write(tempBuf)
	packet.Write(tempBuf)

	buf := new(bytes.Buffer)
	VarInt(tempBuf.Len()).Write(buf)
	buf.Write(tempBuf.Bytes())
	s.Writer.Write(buf.Bytes())
}

func HandleConnection(conn net.Conn) {
	session := Session{conn, conn, conn, Handshake, Serverbound, uuid.New()}
	Sessions[session.uuid.String()] = &session
	for {
		packet, id, len := getPacket(session)
		if packet == nil { //todo
			if session.state != Play {
				session.close()
				break
			}

			if len-1 <= 0 {
				continue
			}

			buf := make([]byte, len-1)
			_, err := io.ReadFull(session.Reader, buf[:])
			if err != nil {
				session.close()
				break
			}
		}

		switch session.state {
		case Handshake:
			if id == 0x00 {
				handshake := *packet.(*HandshakePacket)
				handshake.Read(session)

				if handshake.ProtocolVersion != 47 && handshake.NextState == 2 {
					println("Wrong protocol version: ", handshake.ProtocolVersion)
					session.close()
					break
				}

				if handshake.NextState == 1 {
					session.state = Status
				} else {
					session.state = Login
				}
			}
			break
		case Status:
			if id == 0x00 {
				println("Status")
				status := StatusData{
					Version: VersionStatusData{
						Name:     "123",
						Protocol: 0,
					},
					Players: PlayersStatusData{
						Max:    0,
						Online: 0,
						Sample: []PlayerDataStatus{},
					},
					Description: DescriptionStatusData{
						Text: "AAAAA",
					},
				}

				session.SendPacket(&ServerStatusResponse{Status: status})
			}
			break
		case Login:
			if id == 0x00 {
				login := *packet.(*ClientLoginStart)
				login.Read(session)
				println("Nick: " + login.Name)
				session.SendPacket(&ServerLoginSuccess{
					Uuid: String(session.uuid.String()),
					Name: login.Name,
				})

				session.SendPacket(&ServerJoinGame{
					EntityId:         1,
					Gamemode:         1,
					Dimension:        0,
					Difficulty:       0,
					MaxPlayers:       2,
					LevelType:        "default",
					ReducedDebugInfo: false,
				})

				session.SendPacket(&ServerSpawnPosition{Location: Position{
					X: 0,
					Y: 0,
					Z: 0,
				}})

				session.SendPacket(&ServerPlayerAbilities{
					Invulnerable: false,
					Flying:       true,
					AllowFlying:  true,
					CreativeMode: false,
					FlyingSpeed:  0,
					WalkingSpeed: 0,
				})

				session.SendPacket(&ServerPlayerPosAndLook{
					X:     0,
					Y:     0,
					Z:     0,
					Yaw:   0,
					Pitch: 0,
					Flags: 0,
				})

				session.SendPacket(&ServerChatMessage{
					Message: Msg{
						Text: "§aHello, world!",
					},
					Position: 0,
				})

				session.state = Play
			}

			break
		case Play:
			if id == 0x01 {
				chatMessage := *packet.(*ClientChatMessage)
				chatMessage.Read(session)
				println("Chat: " + chatMessage.Message)
			}
			break
		}
	}
}

func getPacket(session Session) (Packet, VarInt, VarInt) {
	var packetLen VarInt
	packetLen.Read(session)
	if packetLen == 0 {
		return nil, 0, packetLen
	}

	var packetId VarInt
	packetId.Read(session)

	return GetNewPacket(Packets[session.direction][session.state][packetId]), packetId, packetLen
}
