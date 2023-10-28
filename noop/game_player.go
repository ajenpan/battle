package noop

import (
	battle "github.com/ajenpan/battle"
	protobuf "google.golang.org/protobuf/proto"
)

type GamePlayer struct {
	SeatID  int32
	Score   int64
	Role    battle.RoleType
	TableID string
}

func (p *GamePlayer) GetSeatID() int32 {
	return p.SeatID
}

func (p *GamePlayer) GetScore() int64 {
	return p.Score
}

func (p *GamePlayer) SendMessage(protobuf.Message) error {
	return nil
}

func (p *GamePlayer) GetRole() battle.RoleType {
	return p.Role
}
func (p *GamePlayer) GetTableID() string {
	return p.TableID
}
