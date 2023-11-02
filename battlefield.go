package battlefield

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
	BattleStatus_Idle    GameStatus = iota
	BattleStatus_Started GameStatus = iota
	BattleStatus_Over    GameStatus = iota
)

type PlayerMessage struct {
	Head []byte
	Body []byte
}

type PlayerStatus interface {
	IsOnline() bool
	IsJoined() bool
}

type PlayerBattleInfo interface {
	GetSeatID() uint32
	GetScore() int64
}

type Player interface {
	GetUID() uint64
	GetRole() int32

	PlayerStatus
	PlayerBattleInfo
}

type Battle interface {
	GetID() string
	SendPlayerMessage(Player, *PlayerMessage)
	BroadcastPlayerMessage(*PlayerMessage)

	ReportBattleStatus(GameStatus)
	ReportBattleEvent(event proto.Message)
}

type Logic interface {
	OnInit(c Battle, conf interface{}) error
	OnStart() error
	OnTick(time.Duration)
	OnReset()

	OnPlayerMessage(Player, *PlayerMessage)
	OnPlayerStatus(Player)
}