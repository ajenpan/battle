package noop

import (
	"time"

	"google.golang.org/protobuf/proto"

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
func (gl *GameLogic) OnMessage(battle.Player, string, []byte) {

}
func (gl *GameLogic) OnEvent(string, proto.Message) {

}
