package handler

import (
	"fmt"
	"sync"
	"sync/atomic"

	"google.golang.org/protobuf/proto"

	"github.com/ajenpan/surf"

	bf "github.com/ajenpan/battle"
	"github.com/ajenpan/battle/msg"
	"github.com/ajenpan/battle/table"
)

type Handler struct {
	battles        sync.Map
	player2Battles *Player2Battles

	Creator       *bf.LogicCreator
	createCounter uint64
}

func New() *Handler {
	h := &Handler{
		Creator: bf.DefaultLoigcCreator,
	}
	return h
}

func (h *Handler) OnReqStartBattle(s *surf.Context, in *msg.ReqStartBattle) (*msg.RespStartBattle, error) {
	if len(in.PlayerInfos) == 0 {
		return nil, fmt.Errorf("player info is empty")
	}

	logic, err := h.Creator.CreateLogic(in.LogicName, in.LogicVersion)
	if err != nil {
		return nil, err
	}

	battleid := atomic.AddUint64(&h.createCounter, 1)
	if in.BattleConf != nil && in.BattleConf.BattleId > 0 {
		battleid = in.BattleConf.BattleId
		_, ok := h.battles.Load(battleid)
		if ok {
			return nil, fmt.Errorf("battleid already exists")
		}
	}

	closer := func() {
		fmt.Println("battle close :", battleid)
	}

	d := table.NewTable(table.TableOption{
		ID:     battleid,
		Conf:   in.BattleConf,
		Closer: closer,
	})

	players, err := table.NewPlayers(in.PlayerInfos)
	if err != nil {
		return nil, err
	}

	err = d.Init(players, logic, in.LogicConf)
	if err != nil {
		return nil, err
	}

	err = d.Start()
	if err != nil {
		return nil, err
	}

	_, loaded := h.battles.LoadOrStore(battleid, d)
	if loaded {
		d.Close()
		return nil, fmt.Errorf("battleid already exists")
	}

	out := &msg.RespStartBattle{
		BattleId: d.GetID(),
	}
	return out, nil
}

func (h *Handler) OnReqStopBattle(s *surf.Context, in *msg.ReqStopBattle) (*msg.RespStopBattle, error) {
	out := &msg.RespStopBattle{}

	d := h.GetBattleById(in.BattleId)
	if d == nil {
		return out, fmt.Errorf("battle not found")
	}
	d.Close()

	h.battles.Delete(in.BattleId)
	return out, nil
}

func (h *Handler) OnReqJoinBattle(s *surf.Context, in *msg.ReqJoinBattle) (*msg.RespJoinBattle, error) {
	b := h.GetBattleById(in.BattleId)
	if b == nil {
		return nil, fmt.Errorf("battle not found")
	}

	b.OnPlayerJoin(s.UId)

	h.player2Battles.AddPlayerBattle(s.UId, in.BattleId)

	return nil, nil
}

func (h *Handler) OnReqQuitBattle(s *surf.Context, in *msg.ReqQuitBattle) (*msg.RespQuitBattle, error) {
	b := h.GetBattleById(in.BattleId)
	if b == nil {
		return nil, fmt.Errorf("battle not found")
	}
	b.OnPlayerQuit(s.UId)
	h.player2Battles.RemovePlayerBattle(s.UId, in.BattleId)
	return nil, nil
}

func (h *Handler) OnBattleMessageWrap(s *surf.Context, msg *msg.BattleMessageWrap) {
	uid := (s.UId)
	b := h.GetBattleById(msg.BattleId)
	if b == nil {
		return
	}
	b.OnPlayerMessage(uid, &bf.PlayerMessage{
		Head: msg.Head,
		Body: msg.Body,
	})
}

func (h *Handler) GetBattleById(battleId uint64) *table.Table {
	if raw, ok := h.battles.Load(battleId); ok {
		return raw.(*table.Table)
	}
	return nil
}

func (h *Handler) OnEvent(topc string, msg proto.Message) {

}
