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

func (gd *NoopTable) SendMessage(battle.Player, *battle.PlayerMessage) {

}

func (gd *NoopTable) BroadcastMessage(*battle.PlayerMessage) {

}

func (gd *NoopTable) ReportBattleEvent(event proto.Message) {

}

func (gd *NoopTable) ReportBattleStatus(battle.GameStatus) {

}

func (gd *NoopTable) GetTableID() string {
	return ""
}
