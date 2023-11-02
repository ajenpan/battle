package table

import (
	"fmt"

	protobuf "google.golang.org/protobuf/proto"

	bf "github.com/ajenpan/battlefield"
	pb "github.com/ajenpan/battlefield/messages"
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

type playerStatus struct {
	online bool
	joined bool
}

func (ps *playerStatus) IsOnline() bool {
	return ps.online
}
func (ps *playerStatus) IsJoined() bool {
	return ps.joined
}

type Player struct {
	*pb.PlayerInfo
	*playerStatus

	battleid string
	sender   func(msg *bf.PlayerMessage) error
}

func (p *Player) GetScore() int64 {
	return p.PlayerInfo.MainScore
}

func (p *Player) GetBattleID() string {
	return p.battleid
}

func (p *Player) GetUID() uint64 {
	return p.PlayerInfo.Uid
}

func (p *Player) GetSeatID() uint32 {
	return p.PlayerInfo.SeatId
}

func (p *Player) GetRole() int32 {
	return int32(p.PlayerInfo.Role)
}

func (p *Player) Send(msg *bf.PlayerMessage) error {
	return p.sender(msg)
}
