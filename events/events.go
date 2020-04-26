package events

import (
	"context"

	"github.com/sirupsen/logrus"
)

const (
	AppendRecordEvent       EventType = 1
	UpdateRecordUpEvent     EventType = 2
	UpdateRecordDownEvent   EventType = 3
	UpdateRecordStatusEvent EventType = 4
	UpdateDumpStatusEvent   EventType = 5

	GroupRttEvent EventType = 6
)

type EventType int8

type Event struct {
	Typ   EventType
	Ctx   context.Context
	Value interface{}
}

var Bus = make(chan *Event, 64)

var eventMap = make(map[EventType]func(context.Context, interface{}) error)

func RegisterEvent(typ EventType, f func(context.Context, interface{}) error) {
	eventMap[typ] = f
}

func AutoDial(ctx context.Context) error {
	go func() {
		for {
			select {
			case e := <-Bus:
				if f, ok := eventMap[e.Typ]; ok {
					err := f(e.Ctx, e.Value)
					if err != nil {
						logrus.WithField("event-type", e.Typ).WithError(err).WithField("event-value", e.Value).Error("fail to deal with event")
					}
				} else {
					logrus.WithField("event-type", e.Typ).WithField("event-value", e.Value).Error("not found")
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	return nil
}
