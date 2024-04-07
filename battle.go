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

type GameStatus = int32

const (
	GameStatus_Idle     GameStatus = iota
	GameStatus_Starting GameStatus = iota
	GameStatus_Started  GameStatus = iota
	GameStatus_Closing  GameStatus = iota
	GameStatus_Over     GameStatus = iota
)

type PlayerMsg struct {
	Msgid int32
	Body  []byte
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

type User interface {
	UId() uint32
	Role() int32
	Send(proto.Message) error
}

type Player interface {
	GetUID() uint32
	GetRole() int32
	GetSeatID() int32
	GetScore() int64
	GetStatus() PlayerStatusType

	SetStatus(PlayerStatusType)
}

type TallyInfo struct {
}

type Table interface {
	GetID() uint64
	SendPlayerMessage(Player, *PlayerMsg)
	BroadcastPlayerMessage(*PlayerMsg)

	ReportGameStarted()
	ReportGameTally(*TallyInfo)
	ReportGameOver()

	ReportBattleEvent(event proto.Message)

	AfterFunc(func())

	// SendEvent()
}

type Logic interface {
	OnInit(c Table, players []Player, conf interface{}) error

	OnStart()
	OnClose()

	OnTick(time.Duration)

	OnPlayerMessage(Player, *PlayerMsg)
	OnPlayerStatus(Player, PlayerStatusType)
}
