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

type PlayerMsg struct {
	Head []byte
	Body []byte
}

type PlayerStatusType int16

const (
	PlayerStatus_Unjoin  PlayerStatusType = 0x0
	PlayerStatus_Joined  PlayerStatusType = 0x10
	PlayerStatus_Disconn PlayerStatusType = 0x20
	PlayerStatus_Quit    PlayerStatusType = 0x21
)

type PlayerSession interface {
	Send(*PlayerMsg) error
}

type Player interface {
	GetUID() uint32
	GetRole() uint32
	GetSeatID() uint32
	GetScore() int64

	SetStatus(PlayerStatusType)
	GetStatus() PlayerStatusType

	// 必须通过table 来发送消息, 这样 table 可以做一些统一的处理, 比如回放等等
	// Send(*PlayerMsg) error
	SetSender(func(*PlayerMsg) error)
}

type Table interface {
	GetID() uint64
	SendPlayerMessage(Player, *PlayerMsg)
	BroadcastPlayerMessage(*PlayerMsg)

	ReportBattleStatus(GameStatus)
	ReportBattleEvent(event proto.Message)

	AfterFunc(func())

	OnPlayerMessage(uid uint32, m *PlayerMsg)
}

type Logic interface {
	OnInit(c Table, players []Player, conf interface{}) error

	OnTick(time.Duration)
	OnReset()

	OnPlayerMessage(Player, *PlayerMsg)
	OnPlayerStatus(Player, PlayerStatusType)
}
