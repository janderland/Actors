package websocket

import (
	"context"
	"fmt"
	"github.com/janderland/Actors/mailbox"
	"github.com/janderland/Actors/messages"
	"log"
	"sync"
)

const Name = "WebSocket"

type State struct{}

type Peers struct {
	Conn [1]mailbox.PeerMailBox
	Data [1]mailbox.PeerMailBox
}

type Mutators struct {
	Accept func(messages.Accept, State) ([]interface{}, State)
}

type Actor struct {
	State    State
	MailBox  mailbox.MailBox
	Peers    Peers
	Mutators Mutators
}

func (a *Actor) Act(ctx context.Context, group *sync.WaitGroup) {
	defer group.Done()
	for {
		received := a.MailBox.Get(ctx)
		if received == nil {
			a.trace("Stopping")
			return
		}
		toSend, state := a.mutate(received)
		a.State = state
		a.send(toSend)
	}
}

func (a *Actor) mutate(received interface{}) ([]interface{}, State) {
	switch received.(type) {

	case messages.Accept:
		const logFmt = "Got @Accept '%v'"
		a.trace(logFmt, received)

		accept := received.(messages.Accept)
		return a.Mutators.Accept(accept, a.State)

	default:
		const logFmt = "Ignoring %T '%v'"
		a.trace(logFmt, received, received)
	}

	return nil, a.State
}

func (a *Actor) send(toSend []interface{}) {
	for _, message := range toSend {
		switch message.(type) {

		case messages.Conn:
			conn := message.(messages.Conn)
			for _, peer := range a.Peers.Conn {
				peer.Put(conn.Pack())
			}
			const logFmt = "WebSocket: Put @Conn '%v'"
			a.trace(logFmt, message)

		case messages.Data:
			data := message.(messages.Data)
			for _, peer := range a.Peers.Data {
				peer.Put(data.Pack())
			}
			const logFmt = "WebSocket: Put @Data '%v'"
			a.trace(logFmt, message)

		default:
			const logFmt = "WebSocket: Bad Send %T '%v'"
			a.panic(logFmt, message, message)
		}
	}
}

func (a *Actor) trace(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	log.Printf("%s, %s", Name, msg)
}

func (a *Actor) panic(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	panic(fmt.Sprintf("%s: %s", Name, msg))
}
