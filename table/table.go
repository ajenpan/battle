package table

import (
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/ajenpan/battle"
	"github.com/ajenpan/battle/log"
	"github.com/ajenpan/battle/msg"
)

type TableOption struct {
	ID         uint64
	Conf       *msg.BattleConfig
	CloserFunc func()
}

func NewTable(opt *TableOption) *Table {

	if opt == nil {
		opt = &TableOption{}
	}

	ret := &Table{
		TableOption: opt,
		CreateAt:    time.Now(),
		chClosed:    make(chan struct{}),
		players:     make(map[uint32]*Player),
	}

	ret.chAction = make(chan func(), 100)

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

	OverDeadline time.Time

	CreateAt time.Time
	StartAt  time.Time
	OverAt   time.Time
	Age      time.Duration

	logic battle.Logic

	rwlock  sync.RWMutex
	players map[uint32]*Player

	chAction chan func()
	chClosed chan struct{}

	ticker *time.Ticker

	status battle.GameStatus
}

func (d *Table) GetID() uint64 {
	return d.TableOption.ID
}

func (d *Table) Init(players []*Player, logic battle.Logic, logicConf interface{}) error {
	d.rwlock.Lock()
	defer d.rwlock.Unlock()

	battlePlayers := make([]battle.Player, len(players))
	for i, p := range players {
		d.players[p.Uid] = p
		battlePlayers[i] = p
	}

	if logic != nil {
		if err := logic.OnInit(d, battlePlayers, logicConf); err != nil {
			return err
		}
	}

	d.logic = logic

	d.ticker = time.NewTicker(1 * time.Second)

	go func() {
		defer func() {
			d.ticker.Stop()
			d.ticker = nil

		}()

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
			case <-d.chClosed:
				return
			case job, ok := <-d.chAction:
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

func (d *Table) PushAction(f func()) {
	d.chAction <- f
}

func (d *Table) AfterFunc(f func()) {
	d.PushAction(f)
}

func (d *Table) onTick(detle time.Duration) {
	d.Age += detle

	if d.logic != nil {
		d.logic.OnTick(detle)
	}

	switch d.GetStatus() {
	case battle.GameStatus_Idle:
		if d.Age > 10*time.Second {
			d.StartGame()
		}
	case battle.GameStatus_Started:
		if d.Age > time.Duration(d.Conf.MaxBattleTime)*time.Second {
			d.CloseGame()
		}
	default:
	}
}

func (d *Table) close() {
	d.rwlock.Lock()
	defer d.rwlock.Unlock()

	select {
	case <-d.chClosed:
		return
	default:
		close(d.chClosed)
	}

	if d.ticker != nil {
		d.ticker.Stop()
	}

	close(d.chAction)

	if d.CloserFunc != nil {
		d.CloserFunc()
	}
}

func (d *Table) GetStatus() battle.GameStatus {
	return atomic.LoadInt32(&d.status)
}

func (d *Table) UpdateStatus(s battle.GameStatus) bool {
	old := atomic.SwapInt32(&d.status, s)
	return old != s
}

func (d *Table) SendPlayerMessage(p battle.Player, msg *battle.PlayerMsg) {
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

func (d *Table) BroadcastPlayerMessage(msg *battle.PlayerMsg) {
	d.rwlock.RLock()
	defer d.rwlock.RUnlock()

	for _, p := range d.players {
		err := p.Send(msg)
		if err != nil {
			log.Error(err)
		}
	}
}

func (d *Table) IsPlaying() bool {
	return d.status == battle.GameStatus_Started
}

func (d *Table) ReportBattleEvent(event proto.Message) {
	d.PublishEvent(event)
}

func (d *Table) GetPlayer(uid uint32) *Player {
	d.rwlock.RLock()
	defer d.rwlock.RUnlock()
	if p, has := d.players[uid]; has {
		return p
	}
	return nil
}

func (d *Table) GetPlayers() map[uint32]*Player {
	d.rwlock.RLock()
	defer d.rwlock.RUnlock()
	return d.players
}

func (d *Table) OnPlayerDisconn(uid uint32) {
	d.PushAction(func() {
		player := d.GetPlayer(uid)
		if player == nil {
			return
		}
		player.send2client = nil
		d.onPlayerStatusChange(player, battle.PlayerStatus_Disconn)
	})
}

func (d *Table) onPlayerStatusChange(p *Player, currstat battle.PlayerStatusType) {
	p.status = currstat
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

func (d *Table) StartGame() {
	ok := atomic.CompareAndSwapInt32(&d.status, battle.GameStatus_Idle, battle.GameStatus_Starting)
	if !ok {
		return
	}

	if d.logic != nil {
		d.logic.OnStart()
	}

	d.OverDeadline = time.Now().Add(time.Duration(d.TableOption.Conf.MaxBattleTime) * time.Second)
}

func (d *Table) CloseGame() {
	ok := atomic.CompareAndSwapInt32(&d.status, battle.GameStatus_Started, battle.GameStatus_Closing)
	if !ok {
		return
	}
	if d.logic != nil {
		d.logic.OnClose()
	}
}

func (d *Table) ReportGameTally(tally *battle.TallyInfo) {

}

func (d *Table) ReportGameStarted() {
	if ok := d.UpdateStatus(battle.GameStatus_Started); !ok {
		return
	}
	d.StartAt = time.Now()
	d.PublishEvent(&msg.EventBattleStarted{})
}

func (d *Table) ReportGameOver() {
	if ok := d.UpdateStatus(battle.GameStatus_Over); !ok {
		return
	}
	d.OverAt = time.Now()
	d.PublishEvent(&msg.EventBattleOver{})

	d.close()
}

func (d *Table) OnBattleMessageWrap(u battle.User, p *msg.BattleMessageWrap) {
	d.chAction <- func() {
		player := d.GetPlayer(u.UId())
		if player == nil || d.logic == nil {
			return
		}
		switch payload := p.Payload.(type) {
		case *msg.BattleMessageWrap_ReqJoin:
			player.send2client = func(raw *battle.PlayerMsg) error {
				return u.Send(&msg.BattleMessageWrap{
					Battleid: d.ID,
					Payload: &msg.BattleMessageWrap_ToClient{
						ToClient: &msg.MsgToClient{
							Msgid: raw.Msgid,
							Body:  raw.Body,
						}},
				})
			}

			d.onPlayerStatusChange(player, battle.PlayerStatus_Joined)

			if d.GetStatus() == battle.GameStatus_Idle {
				unjoinedCnt := 0
				for _, p := range d.players {
					if p.GetStatus() == battle.PlayerStatus_Unjoin {
						unjoinedCnt++
					}
				}

				if unjoinedCnt == 0 {
					d.StartGame()
				}
			}
			u.Send(&msg.BattleMessageWrap{
				Battleid: d.ID,
				Payload:  &msg.BattleMessageWrap_RespJoin{},
			})
		case *msg.BattleMessageWrap_ReqQuit:
			player.send2client = nil
			d.onPlayerStatusChange(player, battle.PlayerStatus_Quit)

			u.Send(&msg.BattleMessageWrap{
				Battleid: d.ID,
				Payload:  &msg.BattleMessageWrap_RespQuit{},
			})
		case *msg.BattleMessageWrap_ToLogic:
			if d.logic == nil {
				return
			}
			d.logic.OnPlayerMessage(player, &battle.PlayerMsg{
				Msgid: payload.ToLogic.Msgid,
				Body:  payload.ToLogic.Body,
			})
		}
	}
}
