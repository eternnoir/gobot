package gobot

import (
	"github.com/eternnoir/gobot/payload"
)

type Adapter interface {
	Init() error

	Start()

	Reply(orgmessage *payload.Message, text string) error
}
