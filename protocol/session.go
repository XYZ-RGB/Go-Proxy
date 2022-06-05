package protocol

import (
	"bytes"
	"errors"
	"io"
	"net"

	"github.com/google/uuid"
)

var Sessions = make(map[string]*Session)

type Session struct {
	net.Conn
	io.Reader
	io.Writer
	State     ConnectionState
	Direction ConnectionDirection
	Uuid      uuid.UUID
}

func (s Session) close() {
	Sessions[s.Uuid.String()] = nil
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
	Sessions[session.Uuid.String()] = &session
	packetReader := newPacketReader()
	for {
		reader, id, len, err := getPacket(session)
		
		if err != nil {
			session.close()
			break

			// if len-1 <= 0 {
			// 	continue
			// }
			
			// buf := make([]byte, len-1)
			// _, err := io.ReadFull(session.Reader, buf[:])
			// if err != nil {
			// 	session.close()
			// 	break
			// }
		} else {
			packetReader.Send(ReaderData{&session, reader, int32(id), int32(len)})
		}
	}
}

func getPacket(session Session) (io.Reader, VarInt, VarInt, error) {
	var packetLen VarInt
	packetLen.Read(session.Reader)
	if packetLen == 0 {
		return nil, 0, packetLen, errors.New("packetLen is 0")
	}

	var packetId VarInt
	packetIdLen := int(packetId.Read(session.Reader))
	
	buf := make([]byte, (int(packetLen) - packetIdLen))
	io.ReadFull(session.Reader, buf)
	reader := bytes.NewReader(buf)

	return reader, packetId, VarInt(len(buf)), nil
}
