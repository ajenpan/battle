package noop

import (
	battle "github.com/ajenpan/battlefield"
	protobuf "google.golang.org/protobuf/proto"
)

type GamePlayer struct {
	SeatID  uint32
	Score   int64
	Role    battle.RoleType
	TableID string
	UID     uint64
}

func (p *GamePlayer) GetUID() uint64 {
	return p.UID
}

func (p *GamePlayer) GetSeatID() uint32 {
	return p.SeatID
}

func (p *GamePlayer) GetScore() int64 {
	return p.Score
}

func (p *GamePlayer) SendMessage(protobuf.Message) error {
	return nil
}

func (p *GamePlayer) GetRole() int32 {
	return int32(p.Role)
}

func (p *GamePlayer) GetTableID() string {
	return p.TableID
}

func (p *GamePlayer) IsJoined() bool {
	return false
}
func (p *GamePlayer) IsOnline() bool {
	return false
}
