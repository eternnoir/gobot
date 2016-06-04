package gobot

import (
	"github.com/eternnoir/gobot/payload"
)

type Worker interface {
	Init(gobot *Gobot) error
	Process(gobot *Gobot, message *payload.Message) error
}
