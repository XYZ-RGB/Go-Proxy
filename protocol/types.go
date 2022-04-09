package protocol

import (
	"Proxy/chat"
	"io"
	"math"
)

type (
	Boolean       bool
	Byte          byte
	UnsignedByte  uint8
	Short         int16
	UnsignedShort uint16
	Int           int32
	Long          int64
	Float         float32
	Double        float64
	String        string
	VarInt        int32
	VarLong       int64
	Position      struct {
		X, Y, Z int
	}
	Msg chat.Message
)

func (b *Boolean) Read(session Session) {
	var buff Byte
	buff.Read(session)
	*b = Boolean(buff != 0)
}

func (b Boolean) Write(writer io.Writer) {
	var buff Byte
	if b {
		buff = 1
	}
	buff.Write(writer)
}

func (b *Byte) Read(session Session) {
	var buf [1]byte
	_, _ = io.ReadFull(session.Reader, buf[:])
	*b = Byte(buf[0])
}

func (b Byte) Write(writer io.Writer) {
	writer.Write([]byte{byte(b)})
}

func (u *UnsignedByte) Read(session Session) {
	var buf [1]byte
	_, _ = io.ReadFull(session.Reader, buf[:])
	*u = UnsignedByte(buf[0])
}

func (u UnsignedByte) Write(writer io.Writer) {
	writer.Write([]byte{byte(u)})
}

func (s *Short) Read(session Session) {
	var buf [2]byte
	_, _ = io.ReadFull(session.Reader, buf[:])
	*s = Short(int16(buf[0])<<8 | int16(buf[1]))
}

func (s Short) Write(writer io.Writer) {
	writer.Write([]byte{byte(s >> 8), byte(s)})
}

func (u *UnsignedShort) Read(session Session) {
	var buf [2]byte
	_, _ = io.ReadFull(session.Reader, buf[:])
	*u = UnsignedShort(uint16(buf[0])<<8 | uint16(buf[1]))
}

func (u UnsignedShort) Write(writer io.Writer) {
	writer.Write([]byte{byte(u >> 8), byte(u)})
}

func (i *Int) Read(session Session) {
	var buf [4]byte
	_, _ = io.ReadFull(session.Reader, buf[:])
	*i = Int(int32(buf[0])<<24 | int32(buf[1])<<16 | int32(buf[2])<<8 | int32(buf[3]))
}

func (i Int) Write(writer io.Writer) {
	writer.Write([]byte{byte(i >> 24), byte(i >> 16), byte(i >> 8), byte(i)})
}

func (l *Long) Read(session Session) {
	var buf [8]byte
	_, _ = io.ReadFull(session.Reader, buf[:])
	*l = Long(int64(buf[0])<<56 | int64(buf[1])<<48 | int64(buf[2])<<40 | int64(buf[3])<<32 | int64(buf[4])<<24 | int64(buf[5])<<16 | int64(buf[6])<<8 | int64(buf[7]))
}

func (l Long) Write(writer io.Writer) {
	writer.Write([]byte{byte(l >> 56), byte(l >> 48), byte(l >> 40), byte(l >> 32), byte(l >> 24), byte(l >> 16), byte(l >> 8), byte(l)})
}

func (f *Float) Read(session Session) {
	var buf [4]byte
	_, _ = io.ReadFull(session.Reader, buf[:])
	*f = Float(math.Float32frombits(uint32(buf[0])<<24 | uint32(buf[1])<<16 | uint32(buf[2])<<8 | uint32(buf[3])))
}

func (f Float) Write(writer io.Writer) {
	var buf Int
	buf = Int(math.Float32bits(float32(f)))
	buf.Write(writer)
}

func (d *Double) Read(session Session) {
	var buf [8]byte
	_, _ = io.ReadFull(session.Reader, buf[:])
	*d = Double(math.Float64frombits(uint64(buf[0])<<56 | uint64(buf[1])<<48 | uint64(buf[2])<<40 | uint64(buf[3])<<32 |
		uint64(buf[4])<<24 | uint64(buf[5])<<16 | uint64(buf[6])<<8 | uint64(buf[7])))
}

func (d Double) Write(writer io.Writer) {
	var buf Long
	buf = Long(math.Float64bits(float64(d)))
	buf.Write(writer)
}

func (s *String) Read(session Session) {
	var length VarInt
	length.Read(session)
	var buf = make([]byte, length)
	_, _ = io.ReadFull(session.Reader, buf)
	*s = String(buf)
}

func (s String) Write(writer io.Writer) {
	var length VarInt
	length = VarInt(len(s))
	length.Write(writer)
	writer.Write([]byte(s))
}

func (i *VarInt) Read(session Session) {
	var value int32
	var pos int32
	var currentByte byte
	for {
		var buf [1]byte
		_, _ = io.ReadFull(session.Reader, buf[:])
		currentByte = buf[0]

		value |= int32(currentByte&0x7F) << pos
		if (currentByte & 0x80) == 0 {
			break
		}
		pos += 7

		if pos >= 32 {
			break
		}
	}

	*i = VarInt(value)
}

func (i VarInt) Write(writer io.Writer) {
	for {
		var buf [1]byte
		currentByte := byte(i & 0x7F)
		i >>= 7
		if i != 0 {
			currentByte |= 0x80
		}
		buf[0] = currentByte
		writer.Write(buf[:])
		if i == 0 {
			break
		}
	}
}

func (l *VarLong) Read(session Session) {
	var value int64
	var pos int32
	var currentByte byte
	for {
		var buf [1]byte
		_, _ = io.ReadFull(session.Reader, buf[:])
		currentByte = buf[0]

		value |= int64(currentByte&0x7F) << pos
		if (currentByte & 0x80) != 0x80 {
			break
		}
		pos += 7

		if pos >= 64 {
			break
		}
	}
	*l = VarLong(value)
}

func (l VarLong) Write(writer io.Writer) {
	for {
		var buf [1]byte
		currentByte := byte(l & 0x7F)
		l >>= 7
		if l != 0 {
			currentByte |= 0x80
		}
		buf[0] = currentByte
		writer.Write(buf[:])
		if l == 0 {
			break
		}
	}
}

func (p *Position) Read(session Session) {
	var encoded Long
	encoded.Read(session)
	x := int(encoded >> 38)
	y := int(encoded & 0xFFF)
	z := int(encoded << 26 >> 38)

	if x >= 1<<25 {
		x -= 1 << 26
	}
	if y >= 1<<11 {
		y -= 1 << 12
	}
	if z >= 1<<25 {
		z -= 1 << 26
	}

	p.X, p.Y, p.Z = x, y, z

}

func (p Position) Write(writer io.Writer) {
	var b [8]Byte
	position := uint64(p.X&0x3FFFFFF)<<38 | uint64((p.Z&0x3FFFFFF)<<12) | uint64(p.Y&0xFFF)
	for i := 7; i >= 0; i-- {
		b[i] = Byte(position)
		position >>= 8
	}

	for i := 0; i < 8; i++ {
		b[i].Write(writer)
	}
}

func (m *Msg) Read(session Session) {
	var str String
	str.Read(session)

	var message chat.Message
	message.UnmarshalJSON([]byte(str))
	msg := Msg(message)
	m = &msg
}

func (m Msg) Write(writer io.Writer) {
	var message chat.Message
	message = chat.Message(m)
	json, _ := message.MarshalJSON()
	String(json).Write(writer)
}
