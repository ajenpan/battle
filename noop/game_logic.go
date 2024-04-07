package noop

import (
	battle "github.com/ajenpan/battle"
)

func NewGameLogic() battle.Logic {
	return &GameLogic{}
}

type GameLogic struct {
	battle.Logic
}

// func (gl *GameLogic) OnInit(battle.Table, interface{}) error {
// 	return nil
// }

// func (gl *GameLogic) OnPlayerJoin(p []battle.Player) error {
// 	return nil
// }

// func (gl *GameLogic) OnStart() error {
// 	return nil
// }
// func (gl *GameLogic) OnTick( ) {

// }
// func (gl *GameLogic) OnReset() {

// }
// func (gl *GameLogic) OnPlayerMessage(p battle.Player, msg *battle.PlayerMessage) {

// }
// func (gl *GameLogic) OnEvent(topic string, data []byte) {

// }
