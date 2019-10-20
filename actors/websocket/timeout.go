package websocket

import "github.com/janderland/Actors/mailbox"

type _timeout struct {
	handle func(State) ([]interface{}, State)
}

const _timeoutPriority = 1

func (t *_timeout) Pack() mailbox.Message {
	return mailbox.Message{
		Contents: t,
		Priority: _timeoutPriority,
	}
}
