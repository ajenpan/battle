package event

import (
	"github.com/ajenpan/battle/event/proto"
)

type Event = proto.EventMessage

type Publisher interface {
	Publish(e *Event)
}

type Recver interface {
	OnEvent(e *Event)
}
