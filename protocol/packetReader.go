package protocol

import (
	"Proxy/actorSystem"
	"fmt"
	"io"
	"sync/atomic"
)

type counterReader struct {
	reader io.Reader
	count  uint64
}

func (r *counterReader) Read(buf []byte) (int, error) {
    n, err := r.reader.Read(buf)
    atomic.AddUint64(&r.count, uint64(n))
    return n, err
}

func (r *counterReader) Count() uint64 {
    return atomic.LoadUint64(&r.count)
}

type ReaderData struct {
	session *Session
	reader io.Reader
	id int32
	len int32
}

func receivePacket(data ReaderData, actor *actorSystem.Actor[ReaderData]) {
	session := data.session
	id := data.id
	packet := GetNewPacket(Packets[session.Direction][session.State][VarInt(id)])
	if packet != nil {
		reader := counterReader{reader: data.reader}
		packet.Read(&reader)
		session := data.session
		switch session.State {
		case Handshake:
			if id == 0x00 {
				handshake := *packet.(*HandshakePacket)

				if handshake.ProtocolVersion != 47 && handshake.NextState == 2 {
					println("Wrong protocol version: ", handshake.ProtocolVersion)
					session.close()
					break
				}

				if handshake.NextState == 1 {
					session.State = Status
				} else {
					session.State = Login
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
				fmt.Println("Nick: " + login.Name)
				session.SendPacket(&ServerLoginSuccess{
					Uuid: String(session.Uuid.String()),
					Name: login.Name,
				})
				session.State = Play
				actorSystem.NewActorReceiver(actor, 1)

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
			}

			break
		case Play:
			if id == 0x01 {
				chatMessage := *packet.(*ClientChatMessage)
				println("Chat: " + chatMessage.Message)
			}
			break
		}
	}

}

func newPacketReader() actorSystem.Actor[ReaderData] {
	return actorSystem.NewActor(0, 1, receivePacket)
}
