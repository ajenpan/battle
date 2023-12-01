package niuniu

import bf "github.com/ajenpan/battle"

type Robot struct {
	*NNPlayer
}

func (r *Robot) OnMsg(m *bf.PlayerMsg) error {
	return nil
}
