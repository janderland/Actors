package queue

import (
	"context"
	"fmt"
	"github.com/janderland/Actors/mailbox"
	"github.com/janderland/Actors/messages"
	"log"
	"sync"
)

const Name = "Queue"

type State struct{}

type Peers struct {
	Accept [1]mailbox.PeerMailBox
}

type Mutators struct {
	Conn func(messages.Conn, State) ([]interface{}, State)
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

	case messages.Conn:
		const logFmt = "Got @Conn '%v'"
		a.trace(logFmt, received)

		conn := received.(messages.Conn)
		return a.Mutators.Conn(conn, a.State)

	default:
		const logFmt = "Ignoring %T '%v'"
		a.trace(logFmt, received, received)
	}

	return nil, a.State
}

func (a *Actor) send(toSend []interface{}) {
	for _, message := range toSend {
		switch message.(type) {

		case messages.Accept:
			accept := message.(messages.Accept)
			for _, peer := range a.Peers.Accept {
				peer.Put(accept.Pack())
			}

			const logFmt = "Put @Accept '%v'"
			a.trace(logFmt, message)

		default:
			const logFmt = "Bad Send %T '%v'"
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
