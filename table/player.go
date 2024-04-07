package table

import (
	"fmt"

	protobuf "google.golang.org/protobuf/proto"

	"github.com/ajenpan/battle"
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
		// if v.SeatId == 0 {
		// 	return nil, fmt.Errorf("seat id is 0")
		// }
		if v.Uid == 0 {
			return nil, fmt.Errorf("uid is 0")
		}
	}

	return ret, nil
}

type Player struct {
	*pb.PlayerInfo
	battleid    uint64
	send2client func(msg *battle.PlayerMsg) error
	status      battle.PlayerStatusType
}

func (p *Player) GetScore() int64 {
	return p.PlayerInfo.MainScore
}

func (p *Player) GetBattleID() uint64 {
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

func (p *Player) GetSeatID() int32 {
	return p.PlayerInfo.SeatId
}

func (p *Player) GetRole() int32 {
	return int32(p.PlayerInfo.Role)
}

func (p *Player) Send(msg *battle.PlayerMsg) error {
	if p.send2client == nil {
		return fmt.Errorf("sender is nil")
	}
	return p.send2client(msg)
}

func (p *Player) SetSender(f func(*battle.PlayerMsg) error) {
	p.send2client = f
}
