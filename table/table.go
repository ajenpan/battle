package table

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"

	log "github.com/ajenpan/surf/logger"

	bf "github.com/ajenpan/battlefield"
	pb "github.com/ajenpan/battlefield/messages"
)

type TableOption struct {
	ID string
	// EventPublisher event.Publisher
	Conf *pb.TableConfig
	// Closer func(battleid string) error
}

func NewTable(opt TableOption) *Table {
	if opt.ID == "" {
		opt.ID = uuid.NewString()
	}

	ret := &Table{
		TableOption: &opt,
		CreateAt:    time.Now(),
	}

	ret.action = make(chan func(), 100)

	return ret
}

type Table struct {
	*TableOption

	CreateAt time.Time
	StartAt  time.Time
	OverAt   time.Time

	logic bf.Logic

	// watchers    sync.Map
	// evenReport

	rwlock  sync.RWMutex
	players sync.Map

	action chan func()

	ticker *time.Ticker

	battleStatus bf.GameStatus
}

func (d *Table) GetID() string {
	return d.TableOption.ID
}

func (d *Table) Init(players []*Player, logic bf.Logic, logicConf interface{}) error {
	d.rwlock.Lock()
	defer d.rwlock.Unlock()

	if d.logic != nil {
		d.logic.OnReset()
	}

	battlePlayers := make([]bf.Player, len(players))
	for i, p := range players {
		// store player
		d.players.Store(p.Uid, p)

		battlePlayers[i] = p
	}

	if err := logic.OnInit(d, logicConf); err != nil {
		return err
	}

	// if err := logic.OnPlayerJoin(battlePlayers); err != nil {
	// 	return err
	// }

	d.logic = logic

	// switch conf := d.Conf.StartCondition.(type) {
	// case *pb.BattleConfigure_Delayed:
	// 	if conf.Delayed > 0 {
	// 		log.Info("start table after %d seconds", conf.Delayed)
	// 		time.AfterFunc(time.Duration(conf.Delayed)*time.Second, func() {
	// 			err := d.Start()
	// 			if err != nil {
	// 				log.Error(err)
	// 			}
	// 		})
	// 	}
	// }
	return nil
}

func (d *Table) PushAction(f func()) {
	d.action <- f
}

func (d *Table) SetPlayerStatus() {

}

func (d *Table) Start() error {
	d.rwlock.Lock()
	defer d.rwlock.Unlock()

	go func() {
		safecall := func(f func()) {
			defer func() {
				if err := recover(); err != nil {
					log.Errorf("panic: %v", err)
				}
			}()
			f()
		}

		for job := range d.action {
			safecall(job)
		}
	}()

	if d.ticker != nil {
		d.ticker.Stop()
	}

	d.ticker = time.NewTicker(1 * time.Second)
	go func(ticker *time.Ticker) {
		latest := time.Now()
		for now := range ticker.C {
			sub := now.Sub(latest)
			latest = now

			d.PushAction(func() {
				if d.logic != nil {
					d.logic.OnTick(sub)
				}
			})
		}
	}(d.ticker)

	return d.logic.OnStart()
}

func (d *Table) Close() {
	if d.ticker != nil {
		d.ticker.Stop()
	}
	close(d.action)
}

func (d *Table) ReportBattleStatus(s bf.GameStatus) {
	if d.battleStatus == s {
		return
	}

	statusBefore := d.battleStatus
	d.battleStatus = s

	event := &pb.BattleStatusChangeEvent{
		StatusBefore: int32(statusBefore),
		StatusNow:    int32(s),
		BattleId:     string(d.GetID()),
	}
	d.PublishEvent(event)

	switch s {
	case bf.BattleStatus_Idle:
	case bf.BattleStatus_Started:
		d.reportGameStart()
	case bf.BattleStatus_Over:
		d.reportGameOver()
	}
}

func (d *Table) ReportBattleEvent(event proto.Message) {
	d.PublishEvent(event)
}

func (d *Table) SendPlayerMessage(p bf.Player, msg *bf.PlayerMessage) {
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

func (d *Table) BroadcastPlayerMessage(msg *bf.PlayerMessage) {
	d.players.Range(func(key, value interface{}) bool {
		if p, ok := value.(*Player); ok && p != nil {
			err := p.Send(msg)
			if err != nil {
				log.Error(err)
			}
		}
		return true
	})
}

func (d *Table) IsPlaying() bool {
	return d.battleStatus == bf.BattleStatus_Started
}

func (d *Table) reportGameStart() {
	d.StartAt = time.Now()
	d.PublishEvent(&pb.BattleStartEvent{})
}

func (d *Table) reportGameOver() {
	d.OverAt = time.Now()
	d.PublishEvent(&pb.BattleOverEvent{})
}

func (d *Table) GetPlayer(uid uint64) *Player {
	if p, has := d.players.Load(uid); has {
		return p.(*Player)
	}
	return nil
}

func (d *Table) OnPlayerJoin(uid uint64, status int32) {
	d.PushAction(func() {
		player := d.GetPlayer(uid)
		if player == nil {
			return
		}
		player.joined = true
		player.online = true
		d.onPlayerStatusChange(player)
	})
}

func (d *Table) OnPlayerQuit(uid uint64) {
	d.PushAction(func() {
		player := d.GetPlayer(uid)
		if player == nil {
			return
		}
		player.joined = false
		player.online = false
		d.onPlayerStatusChange(player)
	})
}

func (d *Table) OnPlayerDisconn(uid uint64) {
	d.PushAction(func() {
		player := d.GetPlayer(uid)
		if player == nil {
			return
		}
		player.online = false
		d.onPlayerStatusChange(player)
	})
}

func (d *Table) onPlayerStatusChange(p *Player) {
	d.logic.OnPlayerStatus(p)
}

func (d *Table) PublishEvent(event proto.Message) {
	// if d.EventPublisher == nil {
	// 	return
	// }

	// log.Infof("PublishEvent: %s: %v", string(proto.MessageName(event)), event)

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

func (d *Table) OnPlayerMessage(uid uint64, p *bf.PlayerMessage) {
	d.action <- func() {
		player := d.GetPlayer(uid)
		if player != nil && d.logic != nil {
			d.logic.OnPlayerMessage(player, p)
		}
	}
}
