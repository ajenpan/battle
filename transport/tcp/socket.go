package tcp

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync/atomic"
	"time"
)

type SocketStat int32

const (
	Disconnected SocketStat = iota
	Connected    SocketStat = iota
)

type OnMessageFunc func(*Socket, *Packet)
type OnConnStatFunc func(*Socket, SocketStat)
type NewIDFunc func() string

type SocketOptions struct {
	ID string
}

type SocketOption func(*SocketOptions)

var staticIdx uint64

func nextID() string {
	idx := atomic.AddUint64(&staticIdx, 1)
	if idx == 0 {
		idx = atomic.AddUint64(&staticIdx, 1)
	}
	return fmt.Sprintf("tcp_%v_%v", idx, time.Now().Unix())
}

func NewSocket(conn net.Conn, opts SocketOptions) *Socket {
	if opts.ID == "" {
		opts.ID = nextID()
	}

	ret := &Socket{
		id:       opts.ID,
		conn:     conn,
		timeOut:  120 * time.Second,
		chSend:   make(chan *Packet, 10),
		chClosed: make(chan bool),
		state:    Connected,
	}
	return ret
}

type Socket struct {
	conn     net.Conn   // low-level conn fd
	state    SocketStat // current state
	id       string
	chSend   chan *Packet // push message queue
	chClosed chan bool

	timeOut time.Duration

	lastSendAt uint64
	lastRecvAt uint64
}

func (s *Socket) ID() string {
	return s.id
}

func (s *Socket) SendPacket(p *Packet) error {
	if atomic.LoadInt32((*int32)(&s.state)) == int32(Disconnected) {
		return errors.New("sendPacket failed, the socket is disconnected")
	}
	s.chSend <- p
	return nil
}

func (s *Socket) Send(msgid uint32, body []byte) error {
	if len(body) > MaxPacketSize {
		return ErrPacketSizeExcced
	}
	p := &Packet{}
	p.Head.SetType(1)
	p.Head.SetSubType(Packet1SubTypPacket)
	p.Head.SetMsgID(msgid)
	p.Head.SetBodyLength(uint32(len(body)))
	p.Body = body
	return s.SendPacket(p)
}

func (s *Socket) Close() {
	if s == nil {
		return
	}
	stat := atomic.SwapInt32((*int32)(&s.state), int32(Disconnected))
	if stat == int32(Disconnected) {
		return
	}

	if s.conn != nil {
		s.conn.Close()
		s.conn = nil
	}
	close(s.chSend)
	close(s.chClosed)
}

// returns the remote network address.
func (s *Socket) RemoteAddr() net.Addr {
	if s == nil {
		return nil
	}
	return s.conn.RemoteAddr()
}

func (s *Socket) LocalAddr() net.Addr {
	if s == nil {
		return nil
	}
	return s.conn.LocalAddr()
}

// retrun socket work status
func (s *Socket) Status() SocketStat {
	if s == nil {
		return Disconnected
	}
	return SocketStat(atomic.LoadInt32((*int32)(&s.state)))
}

func (s *Socket) writeWork() {
	for p := range s.chSend {
		s.writePacket(p)
	}
}

func writeAll(conn net.Conn, raw []byte) (int, error) {
	writelen := 0
	rawSize := len(raw)

	for writelen < rawSize {
		n, err := conn.Write(raw[writelen:])
		writelen += n
		if err != nil {
			return writelen, err
		}
	}

	return writelen, nil
}

func (s *Socket) readPacket(p *Packet) error {
	if s.Status() == Disconnected {
		return errors.New("recv packet failed, the socket is disconnected")
	}

	var err error

	if s.timeOut > 0 {
		s.conn.SetReadDeadline(time.Now().Add(s.timeOut))
	}

	_, err = io.ReadFull(s.conn, p.Head[:])
	if err != nil {
		return err
	}

	bodylen := p.Head.GetBodyLength()

	if bodylen > 0 {
		//TODO: use buffer pool impove this performance
		p.Body = make([]byte, bodylen)
		_, err = io.ReadFull(s.conn, p.Body)
		return err
	}

	atomic.StoreUint64(&s.lastRecvAt, uint64(time.Now().Unix()))
	return nil
}

func (s *Socket) writePacket(p *Packet) error {
	if s.Status() == Disconnected {
		return errors.New("recv packet failed, the socket is disconnected")
	}

	var err error

	if len(p.Body) >= MaxPacketSize {
		return ErrPacketSizeExcced
	}

	_, err = writeAll(s.conn, p.Head[:])
	if err != nil {
		return err
	}

	if len(p.Body) > 0 {
		_, err = writeAll(s.conn, p.Body)
		if err != nil {
			return err
		}
	}

	atomic.StoreUint64(&s.lastSendAt, uint64(time.Now().Unix()))
	return nil
}
