package handler

import "sync"

type Player2Battles struct {
	rwlock         sync.RWMutex
	player2battles map[uint32]*PlayerBattles
}

func (pb *Player2Battles) GetPlayerBattle(uid uint32) *PlayerBattles {
	pb.rwlock.RLock()
	defer pb.rwlock.RUnlock()
	return pb.player2battles[uid]
}

func (pb *Player2Battles) storePlayerBattle(uid uint32, info *PlayerBattles) *PlayerBattles {
	pb.rwlock.Lock()
	defer pb.rwlock.Unlock()
	if old, has := pb.player2battles[uid]; has {
		return old
	}
	pb.player2battles[uid] = info
	return info
}

func (pb *Player2Battles) deletePlayerBattle(uid uint32) {
	pb.rwlock.Lock()
	defer pb.rwlock.Unlock()
	delete(pb.player2battles, uid)
}

func (pb *Player2Battles) AddPlayerBattle(uid uint32, battleid uint64) {
	infos := pb.GetPlayerBattle(uid)
	if infos != nil {
		infos.JoinBattle(battleid)
	} else {
		infos = &PlayerBattles{}
		infos = pb.storePlayerBattle(uid, infos)
		if infos != nil {
			infos.JoinBattle(battleid)
		}
	}
}

func (pb *Player2Battles) RemovePlayerBattle(uid uint32, battleid uint64) {
	infos := pb.GetPlayerBattle(uid)
	if infos != nil {
		infos.QuitBattle(battleid)
		if infos.Size() == 0 {
			pb.deletePlayerBattle(uid)
		}
	}
}

type PlayerBattleInfo struct {
}

type PlayerBattles struct {
	battles map[uint64]*PlayerBattleInfo
	rwlock  sync.RWMutex
}

func (p *PlayerBattles) Range(f func(battleid uint64, info *PlayerBattleInfo)) {
	p.rwlock.RLock()
	defer p.rwlock.RUnlock()
	for k, v := range p.battles {
		f(k, v)
	}
}
func (p *PlayerBattles) JoinBattle(battleid uint64) {
	p.rwlock.Lock()
	defer p.rwlock.Unlock()

	p.battles[battleid] = &PlayerBattleInfo{}
}

func (p *PlayerBattles) QuitBattle(battleid uint64) {
	p.rwlock.Lock()
	defer p.rwlock.Unlock()

	delete(p.battles, battleid)
}

func (p *PlayerBattles) Size() int {
	p.rwlock.RLock()
	defer p.rwlock.RUnlock()
	return len(p.battles)
}
