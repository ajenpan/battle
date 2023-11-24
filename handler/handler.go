package handler

import (
	"fmt"
	"sync"
	"sync/atomic"

	"google.golang.org/protobuf/proto"

	"route/server"

	bf "github.com/ajenpan/battlefield"
	"github.com/ajenpan/battlefield/msg"
	"github.com/ajenpan/battlefield/table"
	log "github.com/ajenpan/surf/logger"

	"github.com/ajenpan/surf/utils/calltable"
	"github.com/ajenpan/surf/utils/marshal"
)

type Handler struct {
	battles        sync.Map
	player2Battles *Player2Battles

	Creator       *bf.LogicCreator
	CT            *calltable.CallTable[string]
	createCounter uint64
}

type context struct {
	sess server.Session
	msg  *server.Message
}

func New() *Handler {
	h := &Handler{
		Creator: bf.DefaultLoigcCreator,
	}
	ct := calltable.ExtractProtoFile(msg.File_msg_battlefield_proto, h)
	h.CT = ct
	return h
}

func (h *Handler) OnReqStartBattle(s *context, in *msg.ReqStartBattle) (*msg.RespStartBattle, error) {
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

func (h *Handler) OnReqStopBattle(s *context, in *msg.ReqStopBattle) (*msg.RespStopBattle, error) {
	out := &msg.RespStopBattle{}

	d := h.GetBattleById(in.BattleId)
	if d == nil {
		return out, fmt.Errorf("battle not found")
	}
	d.Close()

	h.battles.Delete(in.BattleId)
	return out, nil
}

func (h *Handler) OnReqJoinBattle(s *context, in *msg.ReqJoinBattle) (*msg.RespJoinBattle, error) {
	b := h.GetBattleById(in.BattleId)
	if b == nil {
		return nil, fmt.Errorf("battle not found")
	}

	b.OnPlayerJoin(s.msg.Head.Uid)

	h.player2Battles.AddPlayerBattle(s.msg.Head.Uid, in.BattleId)

	return nil, nil
}

func (h *Handler) OnReqQuitBattle(s *context, in *msg.ReqQuitBattle) (*msg.RespQuitBattle, error) {
	b := h.GetBattleById(in.BattleId)
	if b == nil {
		return nil, fmt.Errorf("battle not found")
	}
	b.OnPlayerQuit(s.msg.Head.Uid)
	h.player2Battles.RemovePlayerBattle(s.msg.Head.Uid, in.BattleId)
	return nil, nil
}

func (h *Handler) OnBattleMessageWrap(s *context, msg *msg.BattleMessageWrap) {
	uid := (s.msg.Head.Uid)
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

func (h *Handler) OnSessionStatus(s server.Session, enable bool) {
	fmt.Println("OnConnect:", s.UserID(), s.SessionID(), enable)
	// if !enable {
	// 	battles := h.player2Battles.GetPlayerBattle(s.UID())
	// 	if battles == nil {
	// 		return
	// 	}
	// 	battles.Range(func(battleid uint64, info *PlayerBattleInfo) {
	// 		b := h.getBattleById(battleid)
	// 		if b == nil {
	// 			return
	// 		}
	// 		b.OnPlayerDisconn(s.UID())
	// 	})
	// 	h.player2Battles.deletePlayerBattle(s.UID())
	// }
}

func (h *Handler) OnTcpMessage(s *server.TcpClient, m *server.Message) {
	var err error

	head := m.Head
	method := h.CT.Get(head.Msgname)
	if method == nil {
		log.Warnf("method not found: %s", head.Msgname)
		return
	}
	pbmarshal := &marshal.ProtoMarshaler{}

	req := method.NewRequest()
	err = pbmarshal.Unmarshal(m.Body, req)
	if err != nil {
		log.Warnf("unmarshal error: %v", err)
		return
	}

	result := method.Call(h, &context{sess: s, msg: m}, req)
	reslen := len(result)

	switch reslen {
	case 1:
		err = result[0].Interface().(error)
	case 2:
		err, _ = result[1].Interface().(error)
	}

	if err != nil {
		log.Warnf("method call error: %v", err)
	}

	if reslen == 2 {
		if m.Head.Msgtype != 1 {
			return
		}
		resp, ok := result[0].Interface().(proto.Message)
		if !ok {
			return
		}
		var resperr *server.Error
		if err != nil {
			resperr = &server.Error{Code: -1, Errmsg: err.Error()}
		}
		s.SendRespMsg(m.Head.Uid, m.Head.Seqid, resp, resperr)
	}
}
