package battle

import (
	"time"

	"google.golang.org/protobuf/proto"
)

type RoleType int

const (
	RoleType_Player = iota
	RoleType_Robot  = iota
)

type GameStatus int16

const (
	BattleStatus_Idle  = 1
	BattleStatus_Start = 1
	BattleStatus_Over  = 1
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
	OnMessage(p Player, topic string, data []byte)
	OnEvent(topic string, event proto.Message)
}
