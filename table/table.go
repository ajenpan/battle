package table

import (
	"fmt"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/ajenpan/battle"
	bf "github.com/ajenpan/battle"
	"github.com/ajenpan/battle/msg"
	"github.com/ajenpan/surf/log"
	"github.com/ajenpan/surf/server"
)

type TableOption struct {
	ID     uint64
	Conf   *msg.BattleConfig
	Closer func()
}

func NewTable(opt TableOption) *Table {
	ret := &Table{
		TableOption: &opt,
		CreateAt:    time.Now(),
		closed:      make(chan struct{}),
		players:     make(map[uint32]*Player),
	}

	ret.action = make(chan func(), 100)

	if ret.Conf == nil {
		ret.Conf = &msg.BattleConfig{}
	}
	if ret.Conf.MaxBattleTime == 0 {
		ret.Conf.MaxBattleTime = 300
	}

	return ret
}

type Table struct {
	*TableOption

	CreateAt time.Time
	StartAt  time.Time
	OverAt   time.Time
	Age      time.Duration

	logic bf.Logic

	// watchers    sync.Map
	// evenReport

	playersRWL sync.RWMutex
	players    map[uint32]*Player

	action chan func()
	closed chan struct{}

	ticker *time.Ticker

	status bf.GameStatus
}

func (d *Table) GetID() uint64 {
	return d.TableOption.ID
}

func (d *Table) Init(players []*Player, logic bf.Logic, logicConf interface{}) error {
	d.playersRWL.Lock()
	defer d.playersRWL.Unlock()

	if d.logic != nil {
		d.logic.OnReset()
	}

	battlePlayers := make([]bf.Player, len(players))
	for i, p := range players {
		// store player
		// d.players.Store(p.Uid, p)

		d.players[p.Uid] = p
		battlePlayers[i] = p
	}

	if err := logic.OnInit(d, battlePlayers, logicConf); err != nil {
		return err
	}

	d.logic = logic
	return nil
}

func (d *Table) PushAction(f func()) {
	d.action <- f
}

func (d *Table) AfterFunc(f func()) {
	d.PushAction(f)
}

func (d *Table) Start() error {
	d.playersRWL.Lock()
	defer d.playersRWL.Unlock()
	if d.logic == nil {
		return fmt.Errorf("logic not init")
	}

	if d.ticker != nil {
		d.ticker.Stop()
	}
	d.ticker = time.NewTicker(1 * time.Second)

	go func() {
		safecall := func(f func()) {
			defer func() {
				if err := recover(); err != nil {
					log.Errorf("panic: %v", err)
				}
			}()
			f()
		}

		latest := time.Now()
		for {
			select {
			case <-d.closed:
				return
			case job, ok := <-d.action:
				if !ok {
					return
				}
				safecall(job)
			case now, ok := <-d.ticker.C:
				if ok {
					sub := now.Sub(latest)
					latest = now
					d.onTick(sub)
				}
			}
		}
	}()
	return nil
}

func (d *Table) onTick(detle time.Duration) {
	d.Age += detle

	d.logic.OnTick(detle)

	switch d.status {
	case bf.GameStatus_Idle:
		if d.Age > 10*time.Second {
			// Do fouce close
			fmt.Println("ready timeout")
		}
	case bf.GameStatus_Started:
		if d.Age > time.Duration(d.Conf.MaxBattleTime)*time.Second {
			// Do fouce close
			fmt.Println("game ready timeout")
			d.Close()
		}
	default:
	}
}

func (d *Table) Close() {
	select {
	case <-d.closed:
		return
	default:
		close(d.closed)
	}

	// if d.logic != nil {
	// 	d.logic.OnReset()
	// }

	if d.ticker != nil {
		d.ticker.Stop()
	}

	close(d.action)

	d.Closer()
}

func (d *Table) ReportBattleStatus(s bf.GameStatus) {
	if d.status == s {
		return
	}

	statusBefore := d.status
	d.status = s

	event := &msg.EventBattleStatusChange{
		StatusBefore: int32(statusBefore),
		StatusNow:    int32(s),
		BattleId:     d.GetID(),
	}
	d.PublishEvent(event)

	switch s {
	case bf.GameStatus_Idle:
	case bf.GameStatus_Started:
		d.reportGameStart()
	case bf.GameStatus_Over:
		d.reportGameOver()
	}
}

func (d *Table) ReportBattleEvent(event proto.Message) {
	d.PublishEvent(event)
}

func (d *Table) SendPlayerMessage(p bf.Player, msg *bf.PlayerMsg) {
	rp, ok := p.(*Player)
	if !ok {
		log.Error("player is not Player")
		return
	}
	err := rp.Send(msg)
	if err != nil {
		log.Error(err)
	}
}

func (d *Table) BroadcastPlayerMessage(msg *bf.PlayerMsg) {
	d.playersRWL.RLock()
	defer d.playersRWL.RUnlock()

	for _, p := range d.players {
		err := p.Send(msg)
		if err != nil {
			log.Error(err)
		}
	}
}

func (d *Table) IsPlaying() bool {
	return d.status == bf.GameStatus_Started
}

func (d *Table) reportGameStart() {
	d.StartAt = time.Now()
	d.PublishEvent(&msg.EventBattleStart{})
}

func (d *Table) reportGameOver() {
	d.OverAt = time.Now()
	d.PublishEvent(&msg.EventBattleOver{})
}

func (d *Table) GetPlayer(uid uint32) *Player {
	d.playersRWL.RLock()
	defer d.playersRWL.RUnlock()
	if p, has := d.players[uid]; has {
		return p
	}
	return nil
}

func (d *Table) OnPlayerJoin(s server.Session, uid uint32) {
	d.PushAction(func() {
		player := d.GetPlayer(uid)
		if player == nil {
			return
		}

		player.sender = func(raw *bf.PlayerMsg) error {
			wrap := &msg.BattleMessageWrap{
				BattleId: d.GetID(),
				Head:     raw.Head,
				Body:     raw.Body,
			}
			return s.SendAsync(player.Uid, wrap)
		}

		player.status = battle.PlayerStatus_Joined
		d.onPlayerStatusChange(player, battle.PlayerStatus_Joined)
	})
}

func (d *Table) OnPlayerQuit(uid uint32) {
	d.PushAction(func() {
		player := d.GetPlayer(uid)
		if player == nil {
			return
		}
		player.status = battle.PlayerStatus_Quit
		player.sender = nil
		d.onPlayerStatusChange(player, 2)
	})
}

func (d *Table) OnPlayerDisconn(uid uint32) {
	d.PushAction(func() {
		player := d.GetPlayer(uid)
		if player == nil {
			return
		}
		player.status = battle.PlayerStatus_Disconn
		d.onPlayerStatusChange(player, 3)
	})
}

func (d *Table) onPlayerStatusChange(p *Player, currstat battle.PlayerStatusType) {
	d.logic.OnPlayerStatus(p, currstat)
}

func (d *Table) PublishEvent(event proto.Message) {
	// if d.EventPublisher == nil {
	// 	return
	// }

	log.Infof("PublishEvent: %s: %v", string(proto.MessageName(event)), event)

	// raw, err := proto.Marshal(event)
	// if err != nil {
	// 	log.Error(err)
	// 	return
	// }
	// warp := &evproto.EventMessage{
	// 	Topic:     string(proto.MessageName(event)),
	// 	Timestamp: time.Now().Unix(),
	// 	Data:      raw,
	// }
	// d.EventPublisher.Publish(warp)
}

func (d *Table) OnPlayerMessage(uid uint32, p *bf.PlayerMsg) {
	d.action <- func() {
		player := d.GetPlayer(uid)
		if player != nil && d.logic != nil {
			d.logic.OnPlayerMessage(player, p)
		}
	}
}
