package gobot

import (
	"github.com/eternnoir/gobot/payload"
)

type Worker interface {
	Init() error
	Process(gobot *Gobot, message *payload.Message) error
}
