package table

import (
	"fmt"

	protobuf "google.golang.org/protobuf/proto"

	bf "github.com/ajenpan/battle"
	pb "github.com/ajenpan/battle/proto"
)

func NewPlayer(p *pb.PlayerInfo) *Player {
	return &Player{
		PlayerInfo: protobuf.Clone(p).(*pb.PlayerInfo),
	}
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
	tableid string

	sender func(msg *bf.PlayerMessage) error
}

func (p *Player) GetScore() int64 {
	return p.PlayerInfo.Score
}

func (p *Player) GetTableID() string {
	return p.tableid
}

func (p *Player) GetUserID() int64 {
	return p.PlayerInfo.Uid
}

func (p *Player) GetSeatID() int32 {
	return p.PlayerInfo.SeatId
}

func (p *Player) GetRole() bf.RoleType {
	if p.PlayerInfo.IsRobot {
		return bf.RoleType_Robot
	}
	return bf.RoleType_Player
}

func (p *Player) Send(msg *bf.PlayerMessage) error {
	return p.sender(msg)
}
