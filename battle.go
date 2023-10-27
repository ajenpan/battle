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
	BattleStatus_Idle    GameStatus = iota
	BattleStatus_Started GameStatus = iota
	BattleStatus_Over    GameStatus = iota
)

type PlayerMessage struct {
	Head []byte
	Body []byte
}

type Player interface {
	SeatID() int32
	Score() int64
	Role() RoleType
}

type PlayerStatus struct {
}

type Table interface {
	SendMessage(Player, proto.Message)
	BroadcastMessage(proto.Message)
	ReportBattleStatus(GameStatus)
	ReportBattleEvent(topic string, event proto.Message)
}

type Logic interface {
	OnInit(c Table, conf interface{}) error
	// OnPlayerConn(p Player, online bool)
	// OnPlayerJoin([]Player) error
	// OnPlayerStatusChange(Player, *PlayerStatus)
	OnStart() error
	OnTick(time.Duration)
	OnReset()
	OnPlayerMessage(p Player, msg *PlayerMessage)
}
