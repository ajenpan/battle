package niuniu

import (
	"fmt"

	bf "github.com/ajenpan/battle"
)

type Robot struct {
	*NNPlayer
}

func (r *Robot) OnMsg(m *bf.PlayerMsg) error {
	fmt.Print("robot on msg: ", m)
	return nil
}
