package commander

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"

	"github.com/ajenpan/battle/event"
	evproto "github.com/ajenpan/battle/event/proto"
	log "github.com/ajenpan/battle/logger"

	bf "github.com/ajenpan/battle"
	pb "github.com/ajenpan/battle/proto"
)

type TableOption struct {
	ID             string
	EventPublisher event.Publisher
	Conf           *pb.CommanderConfigure
}

func NewTable(opt *TableOption) *table {
	if opt.ID == "" {
		opt.ID = uuid.NewString()
	}

	ret := &table{
		TableOption: opt,
		CreateAt:    time.Now(),
	}

	ret.action = make(chan func(), 100)

	return ret
}

type table struct {
	*TableOption

	CreateAt time.Time
	StartAt  time.Time
	OverAt   time.Time

	battle bf.Logic

	IsPlaying bool

	// watchers    sync.Map
	// evenReport

	rwlock  sync.RWMutex
	players sync.Map

	action chan func()

	ticker *time.Ticker

	battleStatus bf.GameStatus
}

func (d *table) Init(players []*Player, logic bf.Logic, logicConf interface{}) error {
	d.rwlock.Lock()
	defer d.rwlock.Unlock()

	if d.battle != nil {
		d.battle.OnReset()
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

	if err := logic.OnPlayerJoin(battlePlayers); err != nil {
		return err
	}

	d.battle = logic

	switch conf := d.Conf.StartCondition.(type) {
	case *pb.CommanderConfigure_Delayed:
		if conf.Delayed > 0 {
			log.Info("start table after %d seconds", conf.Delayed)
			time.AfterFunc(time.Duration(conf.Delayed)*time.Second, func() {
				err := d.Start()
				if err != nil {
					log.Error(err)
				}
			})
		}
	}
	return nil
}
func (d *table) pushAction(f func()) {
	d.action <- f
}

func (d *table) Start() error {
	d.rwlock.Lock()
	defer d.rwlock.Unlock()

	go func(jobque chan func()) {
		for job := range jobque {
			job()
		}

	}(d.action)

	if d.ticker != nil {
		d.ticker.Stop()
	}

	d.ticker = time.NewTicker(1 * time.Second)

	go func(ticker *time.Ticker) {
		latest := time.Now()
		for now := range ticker.C {

			sub := now.Sub(latest)
			latest = now

			d.pushAction(func() {
				if d.battle != nil {
					d.battle.OnTick(sub)
				}
			})
		}
	}(d.ticker)

	return d.battle.OnStart()
}

func (d *table) Close() {

	if d.ticker != nil {
		d.ticker.Stop()
	}
	close(d.action)
}

func (d *table) OnReportBattleStatus(s bf.GameStatus) {
	if d.battleStatus == s {
		return
	}

	statusBefore := d.battleStatus
	d.battleStatus = s

	event := &pb.BattleStatusChangeEvent{
		StatusBefore: int32(statusBefore),
		StatusNow:    int32(s),
		BattleId:     d.ID,
	}
	d.PublishEvent(event)
}

func (d *table) OnReportBattleEvent(topic string, event proto.Message) {
	log.Infof("OnReportBattleEvent: %s: %v", string(proto.MessageName(event)), event)

	//TODO:
	warp := &evproto.EventMessage{
		Topic:     string(proto.MessageName(event)),
		Timestamp: time.Now().Unix(),
	}
	// battle event wrap
	d.PublishEvent(warp)
}

func (d *table) SendMessage(p bf.Player, msg proto.Message) {
	rp := p.(*Player)
	err := rp.SendMessage(msg)
	if err != nil {
		log.Error("send message to player: %v, %s: %v", rp.Uid, string(proto.MessageName(msg)), msg)
	} else {
		log.Debug("send message to player: %v, %s: %v", rp.Uid, string(proto.MessageName(msg)), msg)
	}
}

func (d *table) BroadcastMessage(msg proto.Message) {
	msgname := string(proto.MessageName(msg))
	log.Debugf("BroadcastMessage: %s: %v", msgname, msg)

	raw, err := proto.Marshal(msg)
	if err != nil {
		log.Error(err)
		return
	}

	d.players.Range(func(key, value interface{}) bool {
		if p, ok := value.(*Player); ok && p != nil {
			err := p.Send(msgname, raw)
			if err != nil {
				log.Error(err)
			}
		}
		return true
	})
}

func (d *table) reportGameStart() {
	if d.IsPlaying {
		log.Error("table is playing")
		return
	}
	d.IsPlaying = true
	d.StartAt = time.Now()

	d.PublishEvent(&pb.BattleStartEvent{})
}

func (d *table) reportGameOver() {
	if !d.IsPlaying {
		log.Error("table is not playing")
		return
	}

	d.IsPlaying = false
	d.OverAt = time.Now()

	d.PublishEvent(&pb.BattleOverEvent{})
}

func (d *table) GetPlayer(uid int64) *Player {
	if p, has := d.players.Load(uid); has {
		return p.(*Player)
	}
	return nil
}

func (d *table) PublishEvent(event proto.Message) {
	log.Infof("PublishEvent: %s: %v", string(proto.MessageName(event)), event)

	if d.EventPublisher == nil {
		return
	}
	//TODO:
	warp := &evproto.EventMessage{
		Topic:     string(proto.MessageName(event)),
		Timestamp: time.Now().Unix(),
	}
	d.EventPublisher.Publish(warp)
}

func (d *table) OnPlayerMessage(uid int64, topic string, iraw []byte) {
	// here is not safe
	// msg := proto.Clone(fmsg).(*pb.BattleMessageWrap)

	d.action <- func() {
		p := d.GetPlayer(uid)
		if p != nil && d.battle != nil {
			d.battle.OnMessage(p, topic, iraw)
		}
	}
}
