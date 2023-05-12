package tcp

import (
	"encoding/binary"
	"errors"
)

// Codec constants.
const (
	//16MB
	MaxPacketSize = 1<<(3*8) - 1
	HeadLength    = 4
)

// TODO: rand in packet

var ErrWrongPacketHeadLen = errors.New("wrong packet head len")
var ErrWrongPacketType = errors.New("wrong packet type")
var ErrPacketSizeExcced = errors.New("packet size exceed")
var ErrParseHead = errors.New("parse head error")

// Encode create a packet. packet from the raw bytes slice and then encode to network bytes slice
// -<Type>-|-<SubType>-|-<length>-|-<msgid>-|-<data>-
// -1------|-1---------|-4--------|-4-------|--------

type Packet1SubType = uint8

const (
	Packet1SubTypStatAt_   Packet1SubType = iota
	Packet1SubTypHeartbeat Packet1SubType = iota
	Packet1SubTypAck       Packet1SubType = iota
	Packet1SubTypPing      Packet1SubType = iota
	Packet1SubTypPong      Packet1SubType = iota
	Packet1SubTypPacket    Packet1SubType = iota
	Packet1SubTypError     Packet1SubType = iota
	Packet1SubTypEndAt_    Packet1SubType = iota
)

type PacketHeadRaw [10]byte

// type PacketHead struct {
// 	Typ     uint8
// 	SubTyp  uint8
// 	BodyLen uint32
// 	Msgid   uint32
// }

func (hr *PacketHeadRaw) GetType() uint8 {
	return hr[0]
}
func (hr *PacketHeadRaw) GetSubType() uint8 {
	return hr[1]
}

func (hr *PacketHeadRaw) GetBodyLength() uint32 {
	return binary.LittleEndian.Uint32(hr[2:6])
}

func (hr *PacketHeadRaw) GetMsgID() uint32 {
	return binary.LittleEndian.Uint32(hr[6:10])
}

func (hr *PacketHeadRaw) SetType(t uint8) {
	hr[0] = t
}

func (hr *PacketHeadRaw) SetSubType(t uint8) {
	hr[1] = t
}

func (hr *PacketHeadRaw) SetBodyLength(l uint32) {
	binary.LittleEndian.PutUint32(hr[2:6], l)
}

func (hr *PacketHeadRaw) SetMsgID(id uint32) {
	binary.LittleEndian.PutUint32(hr[6:10], id)
}

// func (hr *PacketHeadRaw) Decode() *PacketHead {
// 	return &PacketHead{
// 		Typ:     hr.GetType(),
// 		SubTyp:  hr.GetSubType(),
// 		BodyLen: hr.GetBodyLength(),
// 		Msgid:   hr.GetMsgID(),
// 	}
// }

// func (hr *PacketHeadRaw) Encode(h *PacketHead) {
// 	hr[0] = h.Typ
// 	hr[1] = h.SubTyp
// 	binary.LittleEndian.PutUint32(hr[2:5], h.BodyLen)
// 	binary.LittleEndian.PutUint32(hr[6:9], h.Msgid)
// }

type Packet struct {
	Head PacketHeadRaw
	Body []byte
}

func (p *Packet) Reset() {
	p.Head = PacketHeadRaw{}
	p.Body = nil
}
