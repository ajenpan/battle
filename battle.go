package battle

import (
	"time"

	"google.golang.org/protobuf/proto"
)

type RoleType int

const (
	RoleType_Player RoleType = iota
	RoleType_Robot  RoleType = iota
)

type GameStatus int16

const (
	BattleStatus_Idle  GameStatus = iota
	BattleStatus_Start GameStatus = iota
	BattleStatus_Over  GameStatus = iota
)

type Player interface {
	SeatID() int32
	Score() int64
	Role() RoleType
}

type Table interface {
	SendMessage(Player, proto.Message)
	BroadcastMessage(proto.Message)

	OnReportBattleStatus(GameStatus)
	OnReportBattleEvent(topic string, event proto.Message)
}

type Logic interface {
	OnInit(c Table, conf interface{}) error
	OnPlayerJoin([]Player) error
	OnStart() error
	OnTick(time.Duration)
	OnReset()
	OnMessage(p Player, msgid int, data []byte)
	OnEvent(topic string, data []byte)
}
