package websocket

import (
	"context"
	"fmt"
	"github.com/janderland/Actors/mailbox"
	"github.com/janderland/Actors/messages"
	"log"
	"net/http"
	"sync"
)

type State struct {
	http.Server
}

type Peers struct {
	Conn [1]mailbox.PeerMailBox
	Data [1]mailbox.PeerMailBox
}

type Handlers struct {
	Accept func(messages.Accept, State) ([]interface{}, State)
}

type Actor struct {
	State    State
	MailBox  mailbox.MailBox
	Peers    Peers
	Handlers Handlers
}

func (a *Actor) Act(ctx context.Context, group *sync.WaitGroup) {
	defer group.Done()
	for {
		received := a.MailBox.Get(ctx)
		if received == nil {
			log.Println("WebSocket: Stopping.")
			return
		}
		switch received.(type) {
		case messages.Accept:
			const logFmt = "WebSocket: Got @Accept '%v'"
			log.Printf(logFmt, received)
			accept := received.(messages.Accept)
			toSend, state := a.Handlers.Accept(accept, a.State)
			a.State = state
			a.send(toSend)
		default:
			const logFmt = "WebSocket: Ignoring %T '%v'"
			log.Printf(logFmt, received, received)
		}
	}
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
			log.Printf(logFmt, message)
		case messages.Data:
			data := message.(messages.Data)
			for _, peer := range a.Peers.Data {
				peer.Put(data.Pack())
			}
			const logFmt = "WebSocket: Put @Data '%v'"
			log.Printf(logFmt, message)
		default:
			const logFmt = "WebSocket: Bad Send %T '%v'"
			panic(fmt.Sprintf(logFmt, message, message))
		}
	}
}
