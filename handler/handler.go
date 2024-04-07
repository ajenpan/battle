package handler

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"google.golang.org/protobuf/proto"

	"github.com/ajenpan/battle"
	"github.com/ajenpan/battle/msg"
	"github.com/ajenpan/battle/table"
)

type Handler struct {
	battles    sync.Map
	Creator    *battle.LogicCreator
	createdIdx uint64

	player2Battles *Player2Battles
}

func New() *Handler {
	h := &Handler{
		Creator: battle.DefaultLoigcCreator,
		player2Battles: &Player2Battles{
			player2battles: make(map[uint32]*PlayerBattles),
		},
	}
	return h
}

func (h *Handler) GetBattleById(battleid uint64) *table.Table {
	if raw, ok := h.battles.Load(battleid); ok {
		return raw.(*table.Table)
	}
	return nil
}

func (h *Handler) OnEvent(topc string, msg proto.Message) {

}

func (h *Handler) ReportBattleStatus() {

}

func (h *Handler) OnUserDiscon(u battle.User) {
	battles := h.player2Battles.GetPlayerBattle(u.UId())

	if battles == nil {
		return
	}

	battles.Range(func(battleid uint64, info *PlayerBattleInfo) {
		t := h.GetBattleById(battleid)
		if t == nil {
			return
		}
		t.OnPlayerDisconn(u.UId())
	})
}

func (h *Handler) OnStartBattle(ctx context.Context, in *msg.ReqStartBattle) (*msg.RespStartBattle, error) {
	if len(in.PlayerInfos) == 0 {
		return nil, fmt.Errorf("player info is empty")
	}

	logic, err := h.Creator.CreateLogic(in.LogicName, in.LogicVersion)
	if err != nil {
		return nil, err
	}

	battleid := atomic.AddUint64(&h.createdIdx, 1)

	opt := &table.TableOption{
		ID:   battleid,
		Conf: in.BattleConf,
	}

	d := table.NewTable(opt)

	players, err := table.NewPlayers(in.PlayerInfos)
	if err != nil {
		return nil, err
	}

	err = d.Init(players, logic, in.LogicConf)
	if err != nil {
		return nil, err
	}

	h.battles.Store(battleid, d)

	out := &msg.RespStartBattle{
		BattleId: d.GetID(),
	}

	for _, p := range in.PlayerInfos {
		h.player2Battles.AddPlayerBattle(p.Uid, battleid)
	}

	opt.CloserFunc = func() {
		raw, ok := h.battles.LoadAndDelete(battleid)
		if !ok {
			return
		}
		d, ok := raw.(*table.Table)
		if !ok {
			return
		}
		pp := d.GetPlayers()
		for _, p := range pp {
			h.player2Battles.RemovePlayerBattle(p.Uid, battleid)
		}
	}
	return out, nil
}

func (h *Handler) OnBattleMessageWrap(u battle.User, msg *msg.BattleMessageWrap) {
	b := h.GetBattleById(msg.Battleid)
	if b == nil {
		return
	}
	b.OnBattleMessageWrap(u, msg)
}
