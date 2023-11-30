package noop

import (
	battle "github.com/ajenpan/battle"
	"google.golang.org/protobuf/proto"
)

func NewGameTable() *NoopTable {
	return &NoopTable{}
}

type NoopTable struct {
}

func (gd *NoopTable) GetID() uint64 {
	return 0
}

func (gd *NoopTable) SendPlayerMessage(battle.Player, *battle.PlayerMsg) {

}

func (gd *NoopTable) BroadcastPlayerMessage(*battle.PlayerMsg) {

}

func (gd *NoopTable) ReportBattleEvent(event proto.Message) {

}

func (gd *NoopTable) ReportBattleStatus(battle.GameStatus) {

}

func (gd *NoopTable) GetTableID() string {
	return ""
}
