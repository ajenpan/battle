package battle

import (
	"time"

	"google.golang.org/protobuf/proto"
)

type RoleType int16

const (
	RoleType_Player RoleType = iota
	RoleType_Robot  RoleType = iota
)

type GameStatus int16

const (
	GameStatus_Idle    GameStatus = iota
	GameStatus_Started GameStatus = iota
	GameStatus_Over    GameStatus = iota
)

type PlayerMessage struct {
	Head []byte
	Body []byte
}

type PlayerStatus interface {
	IsOnline() bool
}

type PlayerBattleInfo interface {
	GetSeatID() uint32
	GetScore() int64
}

type Player interface {
	GetUID() uint32
	GetRole() int32

	PlayerStatus
	PlayerBattleInfo
}

type Table interface {
	GetID() uint64
	SendPlayerMessage(Player, *PlayerMessage)
	BroadcastPlayerMessage(*PlayerMessage)

	ReportBattleStatus(GameStatus)
	ReportBattleEvent(event proto.Message)

	AfterFunc(func())
}

type Logic interface {
	OnInit(c Table, players []Player, conf interface{}) error

	OnTick(time.Duration)
	OnReset()

	OnPlayerMessage(Player, *PlayerMessage)
	OnPlayerStatus(Player)
}
