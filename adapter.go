package gobot

import (
	"github.com/eternnoir/gobot/payload"
)

type Adapter interface {
	Init(gobot *Gobot) error

	Start()

	Send(text string) error

	SendToChat(text, chatroom string) error

	Reply(orgmessage *payload.Message, text string) error
}
