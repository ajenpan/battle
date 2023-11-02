package handler

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/google/uuid"
	protobuf "google.golang.org/protobuf/proto"

	bf "github.com/ajenpan/battlefield"
	proto "github.com/ajenpan/battlefield/messages"
	"github.com/ajenpan/battlefield/table"
	"github.com/ajenpan/surf/tcp"
	"github.com/ajenpan/surf/utils/calltable"
)

type Handler struct {
	battles        sync.Map
	player2Battles *Player2Battles

	Creator       *bf.LogicCreator
	CT            *calltable.CallTable[string]
	createCounter int32
}

func New() *Handler {
	h := &Handler{
		Creator: bf.DefaultLoigcCreator,
	}
	ct := calltable.ExtractProtoFile(proto.File_messages_battle_proto, h)
	h.CT = ct
	return h
}

func (h *Handler) StartBattle(s *tcp.Socket, in *proto.StartBattleRequest) (*proto.StartBattleResponse, error) {
	if len(in.PlayerInfos) == 0 {
		return nil, fmt.Errorf("player info is empty")
	}

	logic, err := h.Creator.CreateLogic(string(in.BattleName))
	if err != nil {
		return nil, err
	}

	atomic.AddInt32(&h.createCounter, 1)

	battleid := uuid.NewString() + fmt.Sprintf("-%d", h.createCounter)
	d := table.NewTable(table.TableOption{
		ID:   battleid,
		Conf: in.TableConf,
		// EventPublisher: h.Publisher,
	})

	players, err := table.NewPlayers(in.PlayerInfos)
	if err != nil {
		return nil, err
	}

	err = d.Init(players, logic, in.BattleConf)
	if err != nil {
		return nil, err
	}

	err = d.Start()
	if err != nil {
		return nil, err
	}

	h.battles.Store(battleid, d)

	out := &proto.StartBattleResponse{
		BattleId: d.GetID(),
	}
	return out, nil
}

func (h *Handler) StopBattle(s *tcp.Socket, in *proto.StopBattleRequest) (*proto.StopBattleResponse, error) {
	out := &proto.StopBattleResponse{}

	d := h.getBattleById(in.BattleId)
	if d == nil {
		return out, fmt.Errorf("battle not found")
	}
	d.Close()

	h.battles.Delete(in.BattleId)
	return out, nil
}

func (h *Handler) PlayerJoinBattle(s *tcp.Socket, in *proto.PlayerJoinBattleRequest) (*proto.PlayerJoinBattleResonse, error) {
	b := h.getBattleById(in.BattleId)
	if b == nil {
		return nil, fmt.Errorf("battle not found")
	}
	b.OnPlayerJoin((s.UId), in.ReadyState)

	h.player2Battles.AddPlayerBattle(s.UID(), in.BattleId)
	return nil, nil
}

func (h *Handler) PlayerQuitBattle(s *tcp.Socket, in *proto.PlayerQuitBattleRequest) (*proto.PlayerQuitBattleResponse, error) {
	b := h.getBattleById(in.BattleId)
	if b == nil {
		return nil, fmt.Errorf("battle not found")
	}
	b.OnPlayerQuit(s.UID())

	h.player2Battles.RemovePlayerBattle(s.UID(), in.BattleId)

	return nil, nil
}

func (h *Handler) OnBattlePlayerMessageWrap(s *tcp.Socket, msg *proto.BattlePlayerMessageWrap) {
	uid := (s.UID())
	b := h.getBattleById(msg.BattleId)
	if b == nil {
		return
	}
	b.OnPlayerMessage(uid, &bf.PlayerMessage{
		Head: msg.Head,
		Body: msg.Body,
	})
}

func (h *Handler) getBattleById(battleId string) *table.Table {
	if raw, ok := h.battles.Load(battleId); ok {
		return raw.(*table.Table)
	}
	return nil
}

func (h *Handler) OnEvent(topc string, msg protobuf.Message) {

}

func (h *Handler) OnMessage(s *tcp.Socket, pkg *tcp.THVPacket) {

}

func (h *Handler) OnConnect(s *tcp.Socket, enable bool) {
	if !enable {
		battles := h.player2Battles.GetPlayerBattle(s.UID())
		if battles == nil {
			return
		}

		battles.Range(func(battleid string, info *PlayerBattleInfo) {
			b := h.getBattleById(battleid)
			if b == nil {
				return
			}
			b.OnPlayerDisconn(s.UID())
		})

		h.player2Battles.DeletePlayerBattle(s.UID())
	}
}
