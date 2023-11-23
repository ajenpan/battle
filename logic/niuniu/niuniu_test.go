package niuniu

import (
	"testing"
	"time"

	battle "github.com/ajenpan/battlefield"
	"github.com/ajenpan/battlefield/noop"
)

type TestTable struct {
	noop.NoopTable
}

func (gd *TestTable) SendPlayerMessage(p battle.Player, msg *battle.PlayerMessage) {

}

func TestOnMessage(t *testing.T) {
	table := &TestTable{}
	ps := []battle.Player{}

	nn := New()
	err := nn.OnInit(table, ps, &Config{Downtime: 10 * time.Second})
	if err != nil {
		t.Error(err)
	}
	player := &noop.GamePlayer{
		SeatID: 1,
		Score:  100,
	}

	nn.addPlayer(player)
	nn.OnPlayerMessage(player, PBMarshal(&GameInfoRequest{}))
}
