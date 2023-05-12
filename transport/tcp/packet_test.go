package tcp

import (
	"math/rand"
	"testing"
)

func TestNewPacketHead(t *testing.T) {
	msgid := rand.Uint32()
	bodylen := rand.Uint32()

	p := &Packet{}
	p.Head.SetType(1)
	p.Head.SetSubType(Packet1SubTypPacket)
	p.Head.SetMsgID(msgid)
	p.Head.SetBodyLength(bodylen)

	if p.Head.GetType() != 1 {
		t.Error("GetType failed")
	}
	if p.Head.GetSubType() != Packet1SubTypPacket {
		t.Error("GetSubType failed")
	}
	if p.Head.GetMsgID() != msgid {
		t.Error("GetMsgID failed")
	}
	if p.Head.GetBodyLength() != bodylen {
		t.Error("GetBodyLength failed")
	}
}
