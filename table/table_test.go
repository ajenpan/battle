package table

import (
	"testing"
	"time"

	"github.com/ajenpan/battle/logic/niuniu"
	"github.com/ajenpan/battle/msg"
)

func newTestTable() *Table {
	MaxBattleTime := int32(5)

	opt := &TableOption{
		ID: 1,
		Conf: &msg.BattleConfig{
			MaxBattleTime: MaxBattleTime,
		},
	}

	d := NewTable(opt)

	players, err := NewPlayers([]*msg.PlayerInfo{
		{
			Uid:       1,
			SeatId:    0,
			MainScore: 1000,
		}, {
			Uid:       2,
			SeatId:    1,
			MainScore: 1000,
		}, {
			Uid:       3,
			SeatId:    2,
			MainScore: 1000,
		}, {
			Uid:       4,
			SeatId:    3,
			MainScore: 1000,
		},
	})
	if err != nil {
		return nil
	}
	err = d.Init(players, niuniu.NewLogic(), []byte("{}"))
	if err != nil {
		return nil
	}
	return d
}

func TestTableCloser(t *testing.T) {

	d := newTestTable()
	if d == nil {
		t.FailNow()
		return
	}
	tid := d.GetID()
	store := make(map[uint64]*Table)
	store[tid] = d

	d.TableOption.CloserFunc = func() {
		delete(store, tid)
	}
	d.StartGame()
	time.Sleep(time.Duration(d.TableOption.Conf.MaxBattleTime+2) * time.Second)

	_, has := store[tid]

	if has {
		t.FailNow()
	}
}
