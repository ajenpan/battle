package handler

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	protobuf "google.golang.org/protobuf/proto"

	"github.com/ajenpan/battle/event"
	"github.com/ajenpan/battle/utils/calltable"
	"github.com/google/uuid"

	battle "github.com/ajenpan/battle"
	"github.com/ajenpan/battle/proto"
	"github.com/ajenpan/battle/table"
)

type Handler struct {
	battles sync.Map

	LogicCreator *battle.GameLogicCreator
	CT           *calltable.CallTable
	Publisher    event.Publisher

	createCounter int32
}

func New() *Handler {
	h := &Handler{
		LogicCreator: &battle.GameLogicCreator{},
	}

	ct := calltable.ExtractProtoFile(proto.File_proto_battle_client_proto, h)
	h.CT = ct

	return h
}

func (h *Handler) CreateBattle(ctx context.Context, in *proto.StartBattleRequest) (*proto.StartBattleResponse, error) {
	if len(in.PlayerInfos) == 0 {
		return nil, fmt.Errorf("player info is empty")
	}

	logic, err := h.LogicCreator.CreateLogic(in.GameName)
	if err != nil {
		return nil, err
	}

	atomic.AddInt32(&h.createCounter, 1)

	battleid := uuid.NewString() + fmt.Sprintf("-%d", h.createCounter)
	d := table.NewTable(table.TableOption{
		ID:             battleid,
		Conf:           in.BattleConf,
		EventPublisher: h.Publisher,
	})

	players, err := table.NewPlayers(in.PlayerInfos)
	if err != nil {
		return nil, err
	}

	err = d.Init(logic, players, in.BattleConf)
	if err != nil {
		return nil, err
	}

	err = d.Start()
	if err != nil {
		return nil, err
	}

	h.battles.Store(battleid, d)

	out := &proto.StartBattleResponse{
		BattleId: d.ID,
	}
	return out, nil
}

func (h *Handler) StopBattle(ctx context.Context, in *proto.StopBattleRequest) (*proto.StopBattleResponse, error) {
	out := &proto.StopBattleResponse{}

	d := h.getBattleById(in.BattleId)
	if d == nil {
		return out, fmt.Errorf("battle not found")
	}
	d.Close()

	h.battles.Delete(in.BattleId)
	return out, nil
}

func (h *Handler) OnEvent(topc string, msg protobuf.Message) {

}

func (h *Handler) OnBattleMessageWrap(uid int64, msg *proto.GameMessageWrap) {
	b := h.getBattleById(msg.BattleId)
	if b == nil {
		return
	}
	b.OnPlayerMessage(uid, int(msg.Msgid), msg.Data)
}

func (h *Handler) getBattleById(battleId string) *table.Table {
	if raw, ok := h.battles.Load(battleId); ok {
		return raw.(*table.Table)
	}
	return nil
}