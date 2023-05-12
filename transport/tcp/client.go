package tcp

import (
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type ClientOption func(*ClientOptions)

type ClientOptions struct {
	RemoteAddress string
	OnMessage     OnMessageFunc
	OnConnStat    OnConnStatFunc
}

func NewClient(opts *ClientOptions) *Client {
	ret := &Client{
		Opt: opts,
	}
	return ret
}

type Client struct {
	*Socket
	Opt   *ClientOptions
	mutex sync.Mutex
}

func (c *Client) Connect() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.Socket != nil {
		c.Socket.Close()
	}

	if c.Opt.RemoteAddress == "" {
		return fmt.Errorf("remote address is empty")
	}

	conn, err := net.DialTimeout("tcp", c.Opt.RemoteAddress, 10*time.Second)
	if err != nil {
		return err
	}

	socket := NewSocket(conn, SocketOptions{})

	//send ack
	p := &Packet{}
	p.Head.SetSubType(Packet1SubTypAck)
	err = socket.writePacket(p)
	if err != nil {
		socket.Close()
		return err
	}

	//read ack
	err = socket.readPacket(p)
	if err != nil {
		socket.Close()
		return err
	}
	if p.Head.GetSubType() != Packet1SubTypAck {
		socket.Close()
		return fmt.Errorf("read ack failed, typ: %d", p.Head.GetSubType())
	}

	c.Socket = socket

	if len(p.Body) > 0 {
		c.id = string(p.Body)
	}

	//here is connect finished
	go func() {
		defer socket.Close()
		go socket.writeWork()

		if c.Opt.OnConnStat != nil {
			c.Opt.OnConnStat(c.Socket, Connected)
			defer c.Opt.OnConnStat(c.Socket, Disconnected)
		}

		go func() {
			tk := time.NewTicker(30 * time.Second)
			defer tk.Stop()

			heartbeatPakcet := &Packet{}
			heartbeatPakcet.Head.SetSubType(Packet1SubTypHeartbeat)

			for {
				select {
				case <-tk.C:
					nowUnix := uint64(time.Now().Unix())
					lastSendAt := atomic.LoadUint64(&socket.lastSendAt)
					if nowUnix-lastSendAt > 30 {
						socket.SendPacket(heartbeatPakcet)
					}
				case <-socket.chClosed:
					return
				}
			}
		}()

		var socketErr error = nil
		for {
			p.Reset()
			if socketErr = socket.readPacket(p); socketErr != nil {
				break
			}
			typ := p.Head.GetSubType()
			switch typ {
			case Packet1SubTypAck:
			case Packet1SubTypHeartbeat:
				continue
			default:
				if c.Opt.OnMessage != nil {
					c.Opt.OnMessage(c.Socket, p)
				}
			}
		}
	}()
	return nil
}
