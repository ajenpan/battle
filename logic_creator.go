package battle

import (
	"fmt"
	"strings"
	"sync"
)

var DefaultLoigcCreator = &LogicCreator{}

func RegisterLogic(name, version string, creator func() Logic) error {
	return DefaultLoigcCreator.Add(strings.Join([]string{name, version}, "-"), creator)
}

type LogicCreator struct {
	Store sync.Map
}

func (c *LogicCreator) Add(name string, creator func() Logic) error {
	c.Store.Store(name, creator)
	return nil
}

func (c *LogicCreator) CreateLogic(name string, version string) (Logic, error) {
	v, has := c.Store.Load(strings.Join([]string{name, version}, "-"))
	if !has {
		return nil, fmt.Errorf("game logic %s not found", name)
	}
	creator := v.(func() Logic)
	return creator(), nil
}
