package table

import (
	"fmt"

	protobuf "google.golang.org/protobuf/proto"

	"github.com/ajenpan/battle"
	bf "github.com/ajenpan/battle"
	pb "github.com/ajenpan/battle/msg"
)

func NewPlayer(p *pb.PlayerInfo) *Player {
	ret := &Player{
		PlayerInfo: protobuf.Clone(p).(*pb.PlayerInfo),
	}
	return ret
}

func NewPlayers(infos []*pb.PlayerInfo) ([]*Player, error) {
	ret := make([]*Player, len(infos))
	for i, info := range infos {
		ret[i] = NewPlayer(info)
	}

	// check seatid
	for _, v := range ret {
		if v.SeatId == 0 {
			return nil, fmt.Errorf("seat id is 0")
		}
		if v.Uid == 0 {
			return nil, fmt.Errorf("uid is 0")
		}
	}

	return ret, nil
}

type Player struct {
	*pb.PlayerInfo

	battleid string
	sender   func(msg *bf.PlayerMsg) error

	status battle.PlayerStatusType
}

func (p *Player) GetScore() int64 {
	return p.PlayerInfo.MainScore
}

func (p *Player) GetBattleID() string {
	return p.battleid
}

func (p *Player) GetUID() uint32 {
	return p.PlayerInfo.Uid
}

func (ps *Player) SetStatus(s battle.PlayerStatusType) {
	ps.status = s
}
func (ps *Player) GetStatus() battle.PlayerStatusType {
	return ps.status
}

func (p *Player) GetSeatID() uint32 {
	return p.PlayerInfo.SeatId
}

func (p *Player) GetRole() uint32 {
	return uint32(p.PlayerInfo.Role)
}

func (p *Player) Send(msg *bf.PlayerMsg) error {
	if p.sender == nil {
		return fmt.Errorf("sender is nil")
	}
	return p.sender(msg)
}

func (p *Player) SetSender(f func(*bf.PlayerMsg) error) {
	p.sender = f
}
