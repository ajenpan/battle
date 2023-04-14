package noop

import (
	"time"

	battle "github.com/ajenpan/battle"
)

func NewGameLogic() battle.Logic {
	return &GameLogic{}
}

type GameLogic struct {
}

func (gl *GameLogic) OnInit(battle.Table, interface{}) error {
	return nil
}
func (gl *GameLogic) OnPlayerJoin(p []battle.Player) error {
	return nil
}
func (gl *GameLogic) OnStart() error {
	return nil
}
func (gl *GameLogic) OnTick(time.Duration) {

}
func (gl *GameLogic) OnReset() {

}
func (gl *GameLogic) OnMessage(p battle.Player, msgid int, data []byte) {

}
func (gl *GameLogic) OnEvent(topic string, data []byte) {

}
