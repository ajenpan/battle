package niuniu

import (
	"encoding/json"
	"fmt"
	"math/rand"
	reflect "reflect"
	"time"

	nncard "github.com/ajenpan/poker_algorithm/niuniu"
	"github.com/ajenpan/surf/log"
	protobuf "google.golang.org/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"

	bf "github.com/ajenpan/battle"

	ct "github.com/ajenpan/surf/utils/calltable"
)

func init() {
	bf.RegisterLogic("niuniu", "1.0.0", NewLogic)

	file_niuniu_proto_init()
	g_ct = extractCallMethod(File_niuniu_proto.Messages(), New())
}

var g_ct *ct.CallTable[string]

func NewLogic() bf.Logic {
	return New()
}

func New() *Niuniu {
	ret := &Niuniu{
		playerMap: make(map[uint16]*NNPlayer),
		info:      &GameInfo{},
		conf:      &Config{},
		log:       log.Default,
		joined:    make(map[uint16]bool),
	}

	ret.stagesDowntime = []int{
		0, 5, 5, 5, 5, 5, 5, 5, 5, 5,
	}

	return ret
}

func extractCallMethod(ms protoreflect.MessageDescriptors, h interface{}) *ct.CallTable[string] {
	const MethodPrefix string = "On"
	refh := reflect.TypeOf(h)

	ret := ct.NewCallTable[string]()
	pbMsgType := reflect.TypeOf((*protobuf.Message)(nil)).Elem()

	for i := 0; i < ms.Len(); i++ {
		msg := ms.Get(i)
		msgName := string(msg.Name())
		method, has := refh.MethodByName(MethodPrefix + msgName)
		if !has {
			continue
		}

		if method.Type.NumIn() != 3 {
			continue
		}

		reqMsgType := method.Type.In(2)
		if reqMsgType.Kind() != reflect.Ptr {
			continue
		}
		if !reqMsgType.Implements(pbMsgType) {
			continue
		}

		m := &ct.Method{
			Func:        method.Func,
			RequestType: reqMsgType.Elem(),
		}
		m.InitPool()
		ret.Add(msgName, m)
	}
	return ret
}

func GetMessageMsgID(msg protoreflect.MessageDescriptor) uint32 {
	MSGIDDesc := msg.Enums().ByName("MSGID")
	if MSGIDDesc == nil {
		return 0
	}
	IDDesc := MSGIDDesc.Values().ByName("ID")
	if IDDesc == nil {
		return 0
	}
	return uint32(IDDesc.Number())
}

type NNPlayer struct {
	raw bf.Player
	*PlayerInfo
	rawHandCards *nncard.NNHandCards
}

type Config struct {
	Downtime int
}

func ParseConfig(raw []byte) (*Config, error) {
	ret := &Config{}
	err := json.Unmarshal(raw, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

type Niuniu struct {
	table bf.Table
	conf  *Config
	log   *log.Logger

	info      *GameInfo
	playerMap map[uint16]*NNPlayer // seatid to player

	gameAge        time.Duration
	stageAge       time.Duration
	stagesDowntime []int

	stageStartAt  int64
	stageDeadline int64

	joined map[uint16]bool //加入人数
}

func PBMarshal(msg protobuf.Message) *bf.PlayerMsg {
	fmt.Println("marshal: ", msg)
	body, _ := protobuf.Marshal(msg)
	return &bf.PlayerMsg{
		Head: []byte(msg.ProtoReflect().Descriptor().FullName().Name()),
		Body: body,
	}
}

func (nn *Niuniu) BroadcastMessage(msg protobuf.Message) {
	nn.table.BroadcastPlayerMessage(PBMarshal(msg))
}

func (nn *Niuniu) Send2Player(p *NNPlayer, msg protobuf.Message) {
	nn.table.SendPlayerMessage(p.raw, PBMarshal(msg))
}

func (nn *Niuniu) OnPlayerStatus(players bf.Player, status bf.PlayerStatusType) {
	p := nn.playerConv(players)
	if p == nil {
		return
	}

	if nn.info.Status == GameStatus_IDLE && status == bf.PlayerStatus_Joined {
		nn.joined[uint16(p.raw.GetSeatID())] = true
	}

}

func (nn *Niuniu) OnInit(d bf.Table, ps []bf.Player, conf interface{}) error {
	switch v := conf.(type) {
	case []byte:
		var err error
		nn.conf, err = ParseConfig(v)
		if err != nil {
			return err
		}
	case *Config:
		nn.conf = v
	default:
		return fmt.Errorf("unknow config type ")
	}
	nn.table = d
	nn.info.Status = GameStatus_IDLE
	nn.gameAge = 0

	for _, v := range ps {
		if v.GetRole() == uint32(bf.RoleType_Robot) {
			nn.addRobot(v)
		} else {
			_, err := nn.addPlayer(v)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (nn *Niuniu) doStart() {
	nn.table.ReportBattleStatus(bf.GameStatus_Started)
	nn.ChangeLogicStep(GameStatus_BEGIN)
}

func (nn *Niuniu) OnCommand(topic string, data []byte) {

}

func (nn *Niuniu) OnPlayerMessage(p bf.Player, pmsg *bf.PlayerMsg) {
	if pmsg == nil {
		return
	}

	mname := string(pmsg.Head)
	m := g_ct.Get(mname)
	if m == nil {
		nn.log.Errorf("can't find method: %s", mname)
		return
	}

	req := m.NewRequest()
	if err := protobuf.Unmarshal(pmsg.Body, req.(protobuf.Message)); err != nil {
		nn.log.Error(err)
		return
	}

	pp := nn.playerConv(p)
	if pp == nil {
		return
	}

	m.Call(nn, pp, req)
}

func (nn *Niuniu) OnGameInfoRequest(pp *NNPlayer, req *GameInfoRequest) {

	resp := &GameInfoResponse{
		GameInfo: nn.info,
	}

	nn.Send2Player(pp, resp)
}

func (nn *Niuniu) checkStat(p *NNPlayer, expect GameStatus) error {
	if nn.getLogicStep() == expect {
		return fmt.Errorf("game status error")
	}
	if p.Status != previousStep(expect) {
		return fmt.Errorf("player status error")
	}
	return nil
}

func (nn *Niuniu) OnPlayerRobBankerReport(nnPlayer *NNPlayer, rep *PlayerRobBankerReport) {
	if err := nn.checkStat(nnPlayer, GameStatus_BANKER); err != nil {
		return
	}
	notice := &NotifyPlayerRobBanker{
		SeatId: (nnPlayer.raw.GetSeatID()),
		Rob:    rep.Rob,
	}
	nnPlayer.BankerRob = rep.Rob
	nn.BroadcastMessage(notice)
}

func (nn *Niuniu) OnPlayerBetRateReport(nnPlayer *NNPlayer, pMsg *PlayerBetRateReport) {
	if err := nn.checkStat(nnPlayer, GameStatus_BET); err != nil {
		return
	}

	nnPlayer.BetRate = pMsg.Rate
	nnPlayer.Status = GameStatus_BET

	notice := &NotifyPlayerBetRate{
		SeatId: nnPlayer.GetSeatId(),
		Rate:   pMsg.Rate,
	}
	nn.BroadcastMessage(notice)
}

func (nn *Niuniu) OnPlayerOutCardRequest(nnPlayer *NNPlayer, pMsg *PlayerOutCardReport) {
	if err := nn.checkStat(nnPlayer, GameStatus_SHOW_CARDS); err != nil {
		return
	}

	nnPlayer.OutCard = &OutCardInfo{
		Cards: nnPlayer.rawHandCards.Bytes(),
		Type:  BullType(nnPlayer.rawHandCards.Type()),
	}
	nnPlayer.Status = GameStatus_SHOW_CARDS

	notice := &NotifyPlayerOutCard{
		SeatId:  (nnPlayer.SeatId),
		OutCard: nnPlayer.OutCard,
	}

	nn.BroadcastMessage(notice)
}

func (nn *Niuniu) addPlayer(p bf.Player) (*NNPlayer, error) {
	ret := &NNPlayer{}
	ret.PlayerInfo = &PlayerInfo{}
	ret.PlayerInfo.SeatId = (p.GetSeatID())
	ret.raw = p
	if _, has := nn.playerMap[uint16(p.GetSeatID())]; has {
		return nil, fmt.Errorf("seat repeat")
	}
	nn.playerMap[uint16(p.GetSeatID())] = ret
	return ret, nil
}

func (nn *Niuniu) addRobot(p bf.Player) (*NNPlayer, error) {
	nnp, err := nn.addPlayer(p)
	if err != nil {
		return nil, err
	}
	robot := &Robot{
		NNPlayer: nnp,
	}
	p.SetSender(robot.OnMsg)

	p.SetStatus(bf.PlayerStatus_Joined)

	// 机器人自动准备
	nn.joined[uint16(p.GetSeatID())] = true

	return nnp, nil
}

func (nn *Niuniu) OnTick(duration time.Duration) {
	nn.gameAge += duration
	nn.stageAge += duration

	switch nn.getLogicStep() {
	case GameStatus_IDLE:
		if len(nn.joined) == len(nn.playerMap) {
			nn.doStart()
		}
	case GameStatus_BEGIN:
		if nn.StepTimeover() {
			nn.NextStep()
		}
	case GameStatus_BANKER:
		if nn.StepTimeover() || nn.checkEndBanker() {
			nn.NextStep()
		}
	case GameStatus_BANKER_NOTIFY:
		if nn.StepTimeover() {
			nn.notifyRobBanker()
			nn.NextStep()
		}
	case GameStatus_BET: // 下注
		if nn.StepTimeover() || nn.checkPlayerStep(GameStatus_BET) {
			nn.NextStep()
		}
	case GameStatus_DEAL_CARDS: // 发牌
		nn.sendCardToPlayer()
		nn.NextStep()
	case GameStatus_SHOW_CARDS: // 开牌
		if nn.StepTimeover() || nn.checkPlayerStep(GameStatus_SHOW_CARDS) {
			nn.NextStep()
		}
	case GameStatus_TALLY:
		nn.beginTally()
		nn.NextStep()
	case GameStatus_OVER:
		if nn.StepTimeover() {
			nn.table.ReportBattleStatus(bf.GameStatus_Over)
		}
	default:
		fmt.Println("unknow step")
		//warn
	}
}

func (nn *Niuniu) OnReset() {

}

func (nn *Niuniu) getLogicStep() GameStatus {
	return nn.info.Status
}

func (nn *Niuniu) getStageDowntime(s GameStatus) time.Duration {
	ret := time.Duration(nn.conf.Downtime) * time.Second
	if int(s) < len(nn.stagesDowntime) {
		ret = time.Duration(nn.stagesDowntime[s]) * time.Second
	}
	return ret
}

func nextStep(status GameStatus) GameStatus {
	nextStep := status + 1
	if nextStep > GameStatus_OVER {
		nextStep = GameStatus_OVER
	}
	return nextStep
}

func previousStep(status GameStatus) GameStatus {
	previousStatus := status - 1
	if previousStatus < GameStatus_IDLE {
		previousStatus = GameStatus_OVER
	}
	return previousStatus
}

func (nn *Niuniu) NextStep() {
	nn.ChangeLogicStep(nextStep(nn.getLogicStep()))
}

func (nn *Niuniu) ChangeLogicStep(s GameStatus) {
	beforeStatus := nn.getLogicStep()
	if s == beforeStatus {
		return
	}
	nn.info.Status = s

	countdown := nn.getStageDowntime(s).Seconds()

	//reset stage time
	nn.stageAge = 0
	nn.stageStartAt = time.Now().Unix()
	nn.stageDeadline = nn.stageStartAt + int64(countdown)

	nn.log.Infof("game step changed, before:%d, now:%d", beforeStatus, s)

	if beforeStatus == s {
		nn.log.Errorf("set same step before:%d, now:%d", beforeStatus, s)
	}

	if beforeStatus != GameStatus_OVER {
		if beforeStatus > s {
			nn.log.Errorf("last step is bigger than now before:%v, now:%v", beforeStatus, s)
		}
	}

	notice := &NotifyGameStatusChange{
		BeforeStatus: beforeStatus,
		CurrStatus:   s,
		CountDown:    int32(countdown),
		StatusAt:     uint32(nn.stageStartAt),
	}

	nn.BroadcastMessage(notice)
	nn.Debug()
}

func (nn *Niuniu) playerConv(p bf.Player) *NNPlayer {
	return nn.getPlayerBySeatId(uint16(p.GetSeatID()))
}

func (nn *Niuniu) getPlayerBySeatId(seatid uint16) *NNPlayer {
	p, ok := nn.playerMap[seatid]
	if ok {
		return p
	}
	return nil
}

func (nn *Niuniu) StepTimeover() bool {
	return nn.stageAge >= nn.getStageDowntime(nn.info.Status)
}

func (nn *Niuniu) checkPlayerStep(expect GameStatus) bool {
	for _, p := range nn.playerMap {
		if p.Status != expect {
			return false
		}
	}
	return true
}

func (nn *Niuniu) checkEndBanker() bool {
	for _, p := range nn.playerMap {
		if p.BankerRob == 0 {
			return false
		}
	}
	return true
}

func (nn *Niuniu) notifyRobBanker() {
	for _, p := range nn.playerMap {
		if p.Status != GameStatus_BANKER {
			p.Status = GameStatus_BANKER
		}
	}

	seats := make([]uint32, 0, len(nn.playerMap))

	var maxRob int32 = -1
	for _, p := range nn.playerMap {
		if (p.BankerRob) > maxRob {
			maxRob = p.BankerRob
			seats = seats[:0]
			seats = append(seats, p.SeatId)
		} else if (p.BankerRob) == maxRob {
			seats = append(seats, p.SeatId)
		}
	}

	if len(seats) == 0 {
		nn.log.Errorf("select bank error maxrob:%d", maxRob)
	}

	index := rand.Intn(len(seats))
	bankSeatId := seats[index]
	banker, ok := nn.playerMap[uint16(bankSeatId)]

	if !ok {
		nn.log.Errorf("banker seatid error. seatid:%d,index:%d", bankSeatId, index)
		return
	}

	banker.IsBanker = true
	//庄家不参与下注.提前设置好状态
	banker.Status = GameStatus_BET

	notice := &NotifyBankerSeat{
		SeatId: bankSeatId,
	}

	nn.BroadcastMessage(notice)
}

func (nn *Niuniu) sendCardToPlayer() {
	deck := nncard.NewNNDeck()
	deck.Shuffle()

	for _, p := range nn.playerMap {
		p.rawHandCards = deck.DealHandCards()
		p.HandCards = p.rawHandCards.Bytes()
		p.Status = GameStatus_DEAL_CARDS
		notice := &NotifyPlayerHandCards{
			SeatId:    p.SeatId,
			HandCards: p.HandCards,
		}
		nn.Send2Player(p, notice)
	}

	for _, p := range nn.playerMap {
		p.rawHandCards.Calculate()
	}
}

func (nn *Niuniu) beginTally() {
	var banker *NNPlayer = nil

	for _, p := range nn.playerMap {
		if p.IsBanker {
			banker = p
			break
		}
	}
	if banker == nil {
		nn.log.Errorf("bank is nil")
		return
	}

	notify := &NotifyGameTally{}
	// notify.TallInfo = make([]*PlayerTallyNotify_TallyInfo, 0)
	// type tally struct {
	// 	UserId int64
	// 	Coins  int32
	// }

	bankerTally := &NotifyGameTally_TallyInfo{
		SeatId: banker.SeatId,
		//Coins:  chips*cardRate*p.BetRate - 100,
	}

	for _, p := range nn.playerMap {
		if p.IsBanker {
			continue
		}
		var chips int32 = 5
		var cardRate int32 = 1

		if banker.rawHandCards.Compare(p.rawHandCards) {
			//底注*倍率*牌型倍率
			cardRate += int32(banker.rawHandCards.Type())
			cardRate = -cardRate
		} else {
			cardRate += int32(p.rawHandCards.Type())
		}
		temp := &NotifyGameTally_TallyInfo{
			SeatId: p.SeatId,
			Coins:  chips * cardRate * p.BetRate,
		}
		// notify.TallInfo = append(notify.TallInfo, temp)
		bankerTally.Coins += temp.Coins
	}

	// notify.TallInfo = append(notify.TallInfo, bankerTally)

	nn.BroadcastMessage(notify)
}

func (nn *Niuniu) resetDesk() {
	nn.playerMap = make(map[uint16]*NNPlayer)
	for _, p := range nn.playerMap {
		p.PlayerInfo.Reset()
		p.PlayerInfo.Status = GameStatus_IDLE
		p.PlayerInfo.SeatId = p.SeatId
	}
	nn.ChangeLogicStep(GameStatus_IDLE)
}

func (nn *Niuniu) Debug() {
	fmt.Println("debug: ", nn.info.String())
}
