package noop

import (
	battle "github.com/ajenpan/battlefield"
	"google.golang.org/protobuf/proto"
)

func NewGameTable() *NoopTable {
	return &NoopTable{}
}

type NoopTable struct {
}

func (gd *NoopTable) GetID() string {
	return ""
}
func (gd *NoopTable) SendPlayerMessage(battle.Player, *battle.PlayerMessage) {

}

func (gd *NoopTable) BroadcastPlayerMessage(*battle.PlayerMessage) {

}

func (gd *NoopTable) ReportBattleEvent(event proto.Message) {

}

func (gd *NoopTable) ReportBattleStatus(battle.BattleStatus) {

}

func (gd *NoopTable) GetTableID() string {
	return ""
}
